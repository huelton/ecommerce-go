package controllers

import (
	"database/sql"
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var orderRepo = repositories.NewOrderRepository(config.DB)

func CreateOrder(c *gin.Context) {
	userID := c.GetInt("user_id")

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := orderRepo.CreateOrder(userID, &order)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock not enough for the product"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to create order"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})
}

func ListOrdersUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	orders, err := orderRepo.ListOrdersUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func OrderPayment(c *gin.Context) {
	userID := c.GetInt("user_id")
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = orderRepo.OrderPayment(userID, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid or not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to update status"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Simulated payment successfully"})
}

func ListAllOrdersAdmin(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should access this route"})
		return
	}
	orders, err := orderRepo.ListAllOrdersAdmin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func CancelOrder(c *gin.Context) {
	userID := c.GetInt("user_id")
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = orderRepo.CancelOrder(userID, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid or canceled"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to cancel order"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}
