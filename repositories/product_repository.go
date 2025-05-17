package repositories

import (
	"database/sql"
	"ecommerce/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) CreateProduct(product *models.Product) error {
	return repo.DB.QueryRow(
		"INSERT INTO produtos(name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id",
		product.Name, product.Description, product.Price, product.Quantity,
	).Scan(&product.ID)
}

func (repo *ProductRepository) ListProducts() ([]models.Product, error) {
	rows, err := repo.DB.Query("SELECT id, name, description, price, quantity FROM produtos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity); err != nil {
			continue
		}
		products = append(products, p)
	}
	return products, nil
}

func (repo *ProductRepository) UpdateProduct(product *models.Product) error {
	_, err := repo.DB.Exec(
		"UPDATE produtos SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5",
		product.Name, product.Description, product.Price, product.Quantity, product.ID,
	)
	return err
}
