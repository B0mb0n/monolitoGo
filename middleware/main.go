package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// Define la estructura exacta del JSON que enviará el frontend cuando el usuario pulse el botón de escanear.

type scanRequest struct {
	Target  string   `json:"target"`
	Modules []string `json:"modules"`
}

var (
	registry *Registry		// Variable global que contendrá el directorio de servicios estáticos y sus estados de salud.
	lb       = NewLoadBalancer()	// Crea la instancia global del balanceador de carga utilizando el algoritmo Round Robin que esta en loadbalancer.
)

// forwardToService manda la petición HTTP a la instancia elegida y regresa el JSON que esa instancia respondió.

func forwardToService(url string, target string) (json.RawMessage, error) {
	body, _ := json.Marshal(map[string]string{"target": target})			// Convierte un mapa con el objetivo (target) en un arreglo de bytes con formato JSON.

	resp, err := http.Post(url+"/scan", "application/json", bytes.NewReader(body))	// Envía una petición POST al servicio, mandando el JSON en el cuerpo.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Asegura que el canal de la respuesta se cierre correctamente al terminar la función para evitar fugas de memoria.

	return io.ReadAll(resp.Body) // Respuesta en el Fronted.
}

// enableCORS permite que el frontend (en otro contenedor/puerto) pueda llamar a este middleware desde el navegador.

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Autoriza únicamente al origen especificado.
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Especifica que los métodos permitidos son POST y OPTIONS.
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

// Si las reglas de CORS se cumplen, permite que la petición continúe hacia la función principal scanHandler.

		if r.Method == http.MethodOptions {
			return
		}
		next(w, r)
	}
}

func scanHandler(w http.ResponseWriter, r *http.Request) {

	// Valida que la ruta de escaneo únicamente acepte peticiones de tipo POST.

	if r.Method != http.MethodPost {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Arreglo donde se irán acumulando las respuestas que devuelvan cada uno de los servicios ejecutados.

	var results []json.RawMessage

	// Recorre uno por uno los módulos que el usuario seleccionó en la interfaz

	for _, module := range req.Modules {

		// 1. Service discovery: ¿qué instancias sanas hay para este módulo?
		instances := registry.HealthyInstances(module)

		// 2. Load balancing: ¿a cuál de ellas le toca esta petición?
		instance, ok := lb.Next(module, instances)
		if !ok {
			http.Error(w, "módulo sin instancias disponibles: "+module, http.StatusServiceUnavailable)
			return
		}

		// Indicando a qué módulo, réplica y URL se le está desviando el trabajo.

		log.Printf("dispatch: %s -> %s (%s)", module, instance.Name, instance.URL)

		// 3. Reenviar la petición real al servicio elegido.
		result, err := forwardToService(instance.URL, req.Target)
		if err != nil {
			http.Error(w, "error llamando a "+module, http.StatusBadGateway)
			return
		}
		results = append(results, result)
	}

	w.Header().Set("Content-Type", "application/json")	// Le avisa al navegador que la respuesta que va a recibir está estructurada como un JSON.
	json.NewEncoder(w).Encode(results)			// Convierte todos los resultados acumulados de los escaneos en formato JSON y los envía de regreso al cliente (el frontend).
}

func main() {
	var err error // Capturar errores de arranque
	registry, err = loadStaticRegistry("services.json")
	if err != nil {
		log.Fatalf("no se pudo cargar services.json: %v", err)
	}

	// Revisa salud de cada instancia cada 5 segundos.
	registry.StartHealthChecks(5 * time.Second)

	// Multiplexor para gestionar las rutas.
	mux := http.NewServeMux()
	mux.HandleFunc("/scan", enableCORS(scanHandler))

	log.Println("middleware escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
