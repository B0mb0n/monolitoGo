package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type scanRequest struct {
	Target string `json:"target"`
}

type scanResult struct {
	Module string `json:"module"`
	Target string `json:"target"`
	Output string `json:"output"`
}

func main() {

	// Conexión a la Base de Datos

	dsn := "host=db port=5432 user=postgres password=postgres dbname=monolito sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// La lista de estrategias a correr.
	strategies := []CheckStrategy{
		HeaderCheckStrategy{},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /scan", func(w http.ResponseWriter, r *http.Request) {
		var req scanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		result := scanResult{
			Module: "http_scan",
			Target: req.Target,
			Output: runStrategies(req.Target, strategies),
		}

		_, err := db.Exec(
			`INSERT INTO scan_results (module, target, output) VALUES ($1, $2, $3)`,
			result.Module, result.Target, result.Output,
		)
		if err != nil {
			log.Printf("error guardando resultado: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("http-scan-service escuchando en :9002")
	log.Fatal(http.ListenAndServe(":9002", mux))
}
