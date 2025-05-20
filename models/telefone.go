package models

type Telefone struct {
	ID          int    `json: "id"`
	UserID      string `json: "user_id"`
	PhoneType   string `json: "phone_type"`
	PhoneNumber string `json: "phone_number"`
}
