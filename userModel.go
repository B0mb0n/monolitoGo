package models

import "database/sql"

// User.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserModel struct {
	DB *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}

// GetAll nos devuelve los usuarios.
func (m *UserModel) GetAll() ([]User, error) {
	rows, err := m.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	return users, nil
}

// GetByID busca el usuario por el ID.
func (m *UserModel) GetByID(id int) (User, error) {
	var u User
	query := "SELECT id, name, email FROM users WHERE id = $1"
	err := m.DB.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}
