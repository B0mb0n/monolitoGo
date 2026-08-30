package controllers

import (
	"encoding/json"
	"net/http"

	"monolitoGo/middleware"
	"monolitoGo/models"
)

// ScanHandler -> POST /scan
func ScanHandler(w http.ResponseWriter, r *http.Request) {
	var req models.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	results, err := middleware.Dispatch(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(results)
}
