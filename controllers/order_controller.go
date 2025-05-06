package controllers

import (
	"database/sql"
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
		var stock int
		err := tx.QueryRow("SELECT price, quantity FROM produtos WHERE id = $1", item.ProductID).Scan(&price, &stock)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Product or not Founded"})
			return
		}

		if item.Quantity > stock {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Stock Enough to the product"})
			return
		}

		//Discount Stock
		_, err = tx.Exec("UPDATE produtos SET quantity = quantity - $1 WHERE id = $2",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro to Update stock"})
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
	rows, err := config.DB.Query("SELECT id, status FROM pedidos WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		o.UserID = userID
		rows.Scan(&o.ID, &o.Status)
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

func OrderPayment(c *gin.Context) {
	userID := c.GetInt("user_id")
	orderID := c.Param("id")

	var status string
	err := config.DB.QueryRow("SELECT status FROM pedidos WHERE id = $1 AND user_id = $2", orderID, userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find order"})
		}
		return
	}

	if status == "pago" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already paid"})
		return
	}

	_, err = config.DB.Exec("UPDATE pedidos SET status = 'pago' WHERE id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro to Update status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mensagem": "Simulated payment with successfully"})
}
