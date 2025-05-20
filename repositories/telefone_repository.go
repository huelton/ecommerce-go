package repositories

import (
	"database/sql"
	"ecommerce/models"
)

type TelefoneRepository struct {
	DB *sql.DB
}

func NewTelefoneRepository(db *sql.DB) *TelefoneRepository {
	return &TelefoneRepository{DB: db}
}

func (repo *TelefoneRepository) CreatePhone(phone *models.Telefone) error {
	var querie = "INSERT INTO enderecos(user_id, phone_type, phone_number) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return repo.DB.QueryRow(
		querie, phone.UserID, phone.PhoneType, phone.PhoneNumber,
	).Scan(&phone.ID)
}

func (repo *TelefoneRepository) GetPhoneByUserId(UserID int) ([]models.Telefone, error) {
	var querie = "SELECT id, user_id, phone_type, phone_number FROM enderecos WEHERE id = $1"
	rows, err := repo.DB.Query(querie, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []models.Telefone
	for rows.Next() {
		var ph models.Telefone
		if err := rows.Scan(&ph.ID, &ph.UserID, &ph.PhoneType, &ph.PhoneNumber); err != nil {
			continue
		}
		phones = append(phones, ph)
	}
	return phones, nil
}
