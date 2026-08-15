package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"monolito-go/controllers"
	"monolito-go/models"
)

func main() {
	// Conexión a la base de datos.
	dsn := "host=db port=5432 user=postgres password=postgres dbname=monolito sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Model y controller.
	userModel := models.NewUserModel(db)
	userController := controllers.NewUserController(userModel)

	// Rutas.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", userController.GetUsers)
	mux.HandleFunc("GET /users/{id}", userController.GetUser)

	log.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", mux)
}
