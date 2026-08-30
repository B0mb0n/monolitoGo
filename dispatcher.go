package middleware

import (
	"fmt"

	"monolitoGo/models"
	"monolitoGo/workers"
)

type scannerFunc func(target string) models.ScanResult

// registry mapea el nombre del módulo al worker correcto.
// Agregar un módulo nuevo aquí.
var registry = map[string]scannerFunc{
	"port_scan": workers.PortScan,
	"http_scan": workers.HTTPScan,
	"tls_scan":  workers.TLSScan,
}

// Dispatch redirige cada módulo solicitado a su worker.
// ToAun sin goroutines - eso en la siguiente tarea del proyecto.
func Dispatch(req models.ScanRequest) ([]models.ScanResult, error) {
	var results []models.ScanResult

	for _, moduleName := range req.Modules {
		scanFunc, exists := registry[moduleName]
		if !exists {
			return nil, fmt.Errorf("módulo desconocido: %s", moduleName)
		}
		results = append(results, scanFunc(req.Target))
	}

	return results, nil
}
