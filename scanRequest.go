package models

// ScanRequest: Lo que el cliente envía, target y qué módulos correr.
type ScanRequest struct {
	Target  string   `json:"target"`
	Modules []string `json:"modules"` // ej: ["port_scan", "http_scan"]
}

// ScanResult: Lo que cada worker devuelve.
type ScanResult struct {
	Module string `json:"module"`
	Target string `json:"target"`
	Output string `json:"output"`
}
