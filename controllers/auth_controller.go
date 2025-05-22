package controllers

import (
	"database/sql"
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Registro de usuário
// @Description Cria um novo usuário no sistema
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Dados do usuário"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid JSON"
// @Failure 500 {object} ErrorResponse ""Error in insert User"
// @Router /register [post]
func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	err := config.DB.QueryRow(
		"INSERT INTO usuarios(name, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, string(hashedPassword), user.IsAdmin,
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in insert User"})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusCreated, gin.H{"token": token})

}

// @Summary Login de usuário
// @Description Autentica um usuário e retorna um token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Credenciais do usuário"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid JSON"
// @Failure 401 {object} ErrorResponse "User not Found or Invalid Password"
// @Failure 500 {object} ErrorResponse "Error in find User"
// @Router /login [post]
func Login(c *gin.Context) {
	var login struct {
		Email    string `json: "email`
		Password string `json: "password`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, name, email, password, is_admin FROM usuarios WHERE email = $1",
		login.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in find User"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusOK, gin.H{"token": token})

}
