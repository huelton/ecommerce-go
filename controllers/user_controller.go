package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var userRepository = repositories.NewUserRepository(config.DB)

func UpdateUser(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should update users"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user.ID = id

	// Fetch the existing user to retain email and password
	existingUser, err := userRepository.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		return
	}

	// Retain email and password
	user.Email = existingUser.Email
	user.Password = existingUser.Password

	err = userRepo.UpdateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(http.StatusOK, user)
}
