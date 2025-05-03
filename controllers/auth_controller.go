package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/huelton/ecommerce-go/config"
	"github.com/huelton/ecommerce-go/models"
	"github.com/huelton/ecommerce-go/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	err := config.DB.QueryRow(
		"INSERT INTO usuarios(nome, email, senha, is_admin) VALUES ($1, $2, $3, $4) RETURNIN id",
		user.Name, user.Email, string(hashedPassword), user.IsAdmin,
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in insert User"})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusCreated, gin.H{"token": token})

}
