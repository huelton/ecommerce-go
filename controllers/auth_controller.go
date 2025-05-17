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

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := userRepo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in insert User"})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.IsAdmin)
	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func Login(c *gin.Context) {
	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	user, err := userRepo.GetUserByEmail(login.Email)
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
