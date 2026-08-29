# Cómputo distribuido

# monolitoGo - Tarea monolito 

# Aarón G. Salto - 0268195

## Archivos y directorios.
```bash
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ tree      
.
├── controllers
│   └── userController.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── models
│   └── userModel.go
└── schema.sql

3 directories, 8 files
```
Cambie de idea de proyecto por lo que la estructura quedo de la siguiente manera. De igual forma dejo lo anterior como evidencia.
```bash
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ tree 
.
├── controllers
│   ├── scanController.go
│   └── userController.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── middleware
│   └── dispatcher.go
├── models
│   ├── scanRequest.go
│   └── userModel.go
├── schema.sql
└── workers
    └── scanners.go
```

## Pruebas.
```bash
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/users  
[{"id":1,"name":"Aarón Salto","email":"aaron@example.com"},{"id":2,"name":"Alejandra","email":"ale@example.com"}]
                                                                                                                                                                                                                                                                                                                            
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/users/1
{"id":1,"name":"Aarón Salto","email":"aaron@example.com"}
                                                                                                                                                                                                                                                                                                                            
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/users/2
{"id":2,"name":"Alejandra","email":"ale@example.com"}
```
Pruebas de nueva idea de proyecto.
```bash
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":["port_scan"]}' 
[{"module":"port_scan","target":"192.168.1.20","output":"Simulado: 22/tcp OPEN, 80/tcp OPEN, 443/tcp OPEN"}]
                                                                                                                                                            
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":["tls_scan"]}' 
[{"module":"tls_scan","target":"192.168.1.20","output":"Simulado: TLS 1.2 soportado, TLS 1.3 soportado, certificado válido"}]
                                                                                                                                                            
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":["http_scan"]}'
[{"module":"http_scan","target":"192.168.1.20","output":"Simulado: Server nginx, falta header X-Content-Type-Options"}]
                                                                                                                                                            
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":["port_scan","tls_scan","http_scan"]}'
[{"module":"port_scan","target":"192.168.1.20","output":"Simulado: 22/tcp OPEN, 80/tcp OPEN, 443/tcp OPEN"},{"module":"tls_scan","target":"192.168.1.20","output":"Simulado: TLS 1.2 soportado, TLS 1.3 soportado, certificado válido"},{"module":"http_scan","target":"192.168.1.20","output":"Simulado: Server nginx, falta header X-Content-Type-Options"}]

┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":["https_scan"]}'                   
módulo desconocido: https_scan
                                                                                                                              
┌──(B0mb0ncito㉿kali)-[~/Documents/monolitoGo]
└─$ curl http://localhost:8080/scan -d '{"target":"192.168.1.20","modules":[]}'  
null
```


