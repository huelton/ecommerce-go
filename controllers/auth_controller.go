package controllers

import (
	"database/sql"
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/repositories"
	"ecommerce/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var userRepo = repositories.NewUserRepository(config.DB)

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
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	err := userRepo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error in insert User"})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusCreated, LoginResponse{Token: token})
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
	var login LoginRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	user, err := userRepo.GetUserByEmail(login.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error in find User"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid Password"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusOK, LoginResponse{Token: token})
}
