package repositories

import (
	"database/sql"
	"ecommerce/models"
)

type ContatoAdicionalRepository struct {
	DB *sql.DB
}

func NewContatoAdicionalRepository(db *sql.DB) *ContatoAdicionalRepository {
	return &ContatoAdicionalRepository{DB: db}
}

func (repo *ContatoAdicionalRepository) CreateAditionalContact(contact *models.ContatoAdicional) error {
	var querie = "INSERT INTO enderecos(user_id, name, email, phone_number, is_spouse) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return repo.DB.QueryRow(
		querie, contact.UserID, contact.Name, contact.Email, contact.PhoneNumber, contact.IsSpouse,
	).Scan(&contact.ID)
}

func (repo *ContatoAdicionalRepository) GetAditionalContactByUserId(UserID int) ([]models.ContatoAdicional, error) {
	var querie = "SELECT id, user_id, name, email, phone_number, is_spouse FROM enderecos WEHERE id = $1"
	rows, err := repo.DB.Query(querie, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.ContatoAdicional
	for rows.Next() {
		var c models.ContatoAdicional
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Email, &c.PhoneNumber, &c.IsSpouse); err != nil {
			continue
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}
