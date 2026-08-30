package models

// ScanRequest es lo que el cliente envía: un target y qué módulos correr.
type ScanRequest struct {
	Target  string   `json:"target"`
	Modules []string `json:"modules"` // ej: ["port_scan", "http_scan"]
}

// ScanResult es lo que cada worker devuelve.
type ScanResult struct {
	Module string `json:"module"`
	Target string `json:"target"`
	Output string `json:"output"`
}
