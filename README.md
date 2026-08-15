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

