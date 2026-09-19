// Carga services.json y corre los health checks dinámicos cada 5s

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

// Instance es una réplica concreta de un servicio.

type Instance struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Registry es el "directorio telefónico" del sistema: sabe qué instancias existen por módulo, y cuáles están vivas ahora mismo.

type Registry struct {
	mu       sync.RWMutex          		// RWMutex permite múltiples lectores simultáneos, pero bloquea todo cuando alguien escribe, más lecturas que escrituras, que es este caso.
	all      map[string][]Instance 		// Mapa que almacena todos los servicios estáticos leídos desde el archivo JSON, agrupados por módulo.
	healthy  map[string][]Instance 		// subconjunto que respondió su /health recientemente
}

// loadStaticRegistry lee services.json — service discovery sencillo con documento.

func loadStaticRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var all map[string][]Instance
	if err := json.Unmarshal(data, &all); err != nil {
		return nil, err
	}

	return &Registry{
		all:     all,
		healthy: map[string][]Instance{},
	}, nil
}

// StartHealthChecks lanza una goroutine que cada N segundos golpea GET /health de cada instancia registrada.
// El registro estático dice QUIÉN existe, pero la parte dinámica decide QUIÉN está disponible en este momento.

func (r *Registry) StartHealthChecks(interval time.Duration) {
	check := func() {
		client := http.Client{Timeout: 2 * time.Second} // Evitar que una instancia caída bloquee el proceso indefinidamente.
		newHealthy := map[string][]Instance{}

		for module, instances := range r.all {
			for _, inst := range instances {
				resp, err := client.Get(inst.URL + "/health")
				if err == nil && resp.StatusCode == http.StatusOK {
					newHealthy[module] = append(newHealthy[module], inst)
				}
			}
		}

		// Bloquea de forma segura el registro antes de reemplazar el mapa antiguo de instancias sanas por el nuevo (newHealthy), evitando race conditions.

		r.mu.Lock()
		r.healthy = newHealthy
		r.mu.Unlock()
	}

	// Ejecuta una verificación inmediata al arrancar y luego lanza una goroutine en segundo plano que repetirá el proceso cíclicamente según el interval.

	check() 
	go func() {
		for range time.Tick(interval) {
			check()
		}
	}()
}

// HealthyInstances devuelve las réplicas vivas de un módulo.

func (r *Registry) HealthyInstances(module string) []Instance {
	r.mu.RLock()			// Adquiere un candado de sólo lectura (RLock), que permite consultar instancias sanas al mismo tiempo sin bloquearse entre sí, pero esperando si hay una actualización en curso.
	defer r.mu.RUnlock()		// Asegura que el candado de lectura se libere automáticamente justo antes de que la función termine de retornar.
	return r.healthy[module]	// Retorna la lista actual de réplicas operativas para el módulo consultado.
}
