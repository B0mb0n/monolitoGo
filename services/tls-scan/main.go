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
	dsn := "host=db port=5432 user=postgres password=postgres dbname=monolito sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

		// La fábrica construye cada checker por nombre.

		checks := []string{"cert_validity", "protocol_version"}
		output := ""
		for i, kind := range checks {
			checker := NewChecker(kind)
			if i > 0 {
				output += " | "
			}
			output += checker.Run(req.Target)
		}

		result := scanResult{Module: "tls_scan", Target: req.Target, Output: output}

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

	log.Println("tls-scan-service escuchando en :9003")
	log.Fatal(http.ListenAndServe(":9003", mux))
}
