package models

type User struct {
	ID       int    `json: "id"`
	Name     string `json: "name"`
	Email    string `json: "email"`
	Password string `json: "password"`
	IsAdmin  bool   `json: "is_admin"`
}
