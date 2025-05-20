package repositories

import (
	"database/sql"
	"ecommerce/models"
)

type EnderecoRepository struct {
	DB *sql.DB
}

func NewEnderecoRepository(db *sql.DB) *EnderecoRepository {
	return &EnderecoRepository{DB: db}
}

func (repo *EnderecoRepository) CreateEndereco(address *models.Endereco) error {
	var querie = "INSERT INTO enderecos(user_id, address_type, street_type, street_name, address_complement, address_number, zipcode, neighborhood, address_city, address_state) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	return repo.DB.QueryRow(
		querie, address.UserID, address.AddressType, address.StreetType, address.StreetName, address.AddressComplement, address.Number, address.Zipcode, address.Neighborhood, address.AddressCity, address.AddressState,
	).Scan(&address.ID)
}

func (repo *EnderecoRepository) GetEnderecosByUserId(UserID int) ([]models.Endereco, error) {
	var querie = "SELECT id, user_id, address_type, street_type, street_name, address_complement, address_number, zipcode, neighborhood, address_city, address_state, is_active FROM enderecos WEHERE id = $1"
	rows, err := repo.DB.Query(querie, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Endereco
	for rows.Next() {
		var e models.Endereco
		if err := rows.Scan(&e.ID, &e.UserID, &e.AddressType, &e.StreetType, &e.StreetName, &e.AddressComplement, &e.Number, &e.Zipcode, &e.Neighborhood, &e.AddressCity, &e.AddressState); err != nil {
			continue
		}
		addresses = append(addresses, e)
	}
	return addresses, nil
}
