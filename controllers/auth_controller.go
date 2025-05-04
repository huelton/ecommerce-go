package controllers

import (
	"database/sql"
	"ecommerce-go/config"
	"ecommerce-go/models"
	"ecommerce-go/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

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
