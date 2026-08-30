package workers

import "monolitoGo/models"

func PortScan(target string) models.ScanResult {
	return models.ScanResult{
		Module: "port_scan", Target: target,
		Output: "Simulado: 22/tcp OPEN, 80/tcp OPEN, 443/tcp OPEN",
	}
}

func HTTPScan(target string) models.ScanResult {
	return models.ScanResult{
		Module: "http_scan", Target: target,
		Output: "Simulado: Server nginx, falta header X-Content-Type-Options",
	}
}

func TLSScan(target string) models.ScanResult {
	return models.ScanResult{
		Module: "tls_scan", Target: target,
		Output: "Simulado: TLS 1.2 soportado, TLS 1.3 soportado, certificado válido",
	}
}
