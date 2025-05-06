package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateOrder(c *gin.Context) {
	userID := c.GetInt("user_id")

	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to inicialize Transaction"})
		return
	}

	err = tx.QueryRow("INSERT INTO pedidos(user_id) VALUES ($1) RETURNING id", userID).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to create order"})
		return
	}

	for _, item := range order.Items {
		var price float64
		err := tx.QueryRow("SELECT price FROM produtos WHERE id = $1", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Product"})
			return
		}

		_, err = tx.Exec("INSERT INTO itens_pedido(order_id, product_id, quantity, unit_price) VALUES ($1,$2,$3,$4)",
			order.ID, item.ProductID, item.Quantity, price,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro to insert Items"})
			return
		}
	}

	// Commit da transação para salvar os dados
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Order created successfully"})
}

func ListOrdersUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := config.DB.Query("SELECT id FROM pedidos WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		o.UserID = userID
		rows.Scan(&o.ID)
		items, _ := config.DB.Query("SELECT id, order_id, product_id, quantity, unit_price FROM itens_pedido WHERE order_id = $1", o.ID)
		for items.Next() {
			var item models.OrderItems
			items.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice)
			o.Items = append(o.Items, item)
		}

		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}
