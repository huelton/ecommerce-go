package repositories

import (
	"database/sql"
	"ecommerce/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(user *models.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	return repo.DB.QueryRow(
		"INSERT INTO usuarios(name, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, string(hashedPassword), user.IsAdmin,
	).Scan(&user.ID)
}

func (repo *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.DB.QueryRow(
		"SELECT id, name, email, password, is_admin FROM usuarios WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) UpdateUser(user *models.User) error {
	_, err := repo.DB.Exec(
		"UPDATE usuarios SET name = $1, is_admin = $2 WHERE id = $3",
		user.Name, user.IsAdmin, user.ID,
	)
	return err
}
