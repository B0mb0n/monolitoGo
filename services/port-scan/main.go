package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// scanRequest es lo que este servicio recibe del middleware.

type scanRequest struct {
	Target string `json:"target"`
}

func main() {

	// Conexión a la base de datos compartida.

	dsn := "host=db port=5432 user=postgres password=postgres dbname=monolito sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := NewPostgresRepository(db)

	mux := http.NewServeMux()

	// El middleware usa esto para saber si esta instancia sigue viva (es la base del "service discovery dinámico").

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Endpoint que el middleware llama para pedir el escaneo.

	mux.HandleFunc("POST /scan", func(w http.ResponseWriter, r *http.Request) {
		var req scanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Escaneo simulado.

		result := ScanResult{
			Module: "port_scan",
			Target: req.Target,
			Output: "Simulado: 22/tcp OPEN, 80/tcp OPEN, 443/tcp OPEN",
		}

		if err := repo.Save(result); err != nil {
			log.Printf("error guardando resultado: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("port-scan-service escuchando en :9001")
	log.Fatal(http.ListenAndServe(":9001", mux))
}
