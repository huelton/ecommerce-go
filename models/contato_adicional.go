package models

type ContatoAdicional struct {
	ID          int    `json: "id"`
	UserID      string `json: "user_id"`
	Name        string `json: "name"`
	Email       string `json: "email"`
	PhoneNumber string `json: "phone_number"`
	IsSpouse    bool   `json: "is_spouse"`
}
