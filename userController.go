package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"monolitoGo/models"
)

type UserController struct {
	Model *models.UserModel
}

func NewUserController(m *models.UserModel) *UserController {
	return &UserController{Model: m}
}

// GET /users
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.Model.GetAll()
	if err != nil {
		http.Error(w, "error obteniendo usuarios", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	user, err := c.Model.GetByID(id)
	if err != nil {
		http.Error(w, "usuario no encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}
