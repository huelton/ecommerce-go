package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuccessResponse struct {
	Success string `json:"success"`
}

// CreateOrder godoc
// @Summary Cria um novo pedido
// @Description Cria um pedido para o usuário autenticado com os itens fornecidos
// @Tags Pedidos
// @Accept json
// @Produce json
// @Param order body models.Order true "Dados do pedido"
// @Success 201 {object} SuccessResponse "Order created successfully"
// @Failure 400 {object} ErrorResponse "Invalid JSON"
// @Failure 400 {object} ErrorResponse "Invalid Product"
// @Failure 500 {object} ErrorResponse "Erro to start transaction"
// @Failure 500 {object} ErrorResponse "Error to create order"
// @Failure 500 {object} ErrorResponse "Error to insert Item Order"
// @Router /orders [post]
// @Security BearerAuth
func CreateOrder(c *gin.Context) {
	userID := c.GetInt("user_id")
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON"})
		return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Erro to start transaction"})
		return
	}

	err = tx.QueryRow("INSERT INTO pedidos (user_id) VALUES ($1) RETURNING id", userID).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to create order"})
		return
	}

	for _, item := range order.Items {

		var price float64
		err := tx.QueryRow("SELECT price FROM produtos WHERE id = $1", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Product"})
			return
		}
		_, err = tx.Exec("INSERT INTO items_pedido (order_id, product_id, quantity, unit_price) VALUES ($1,$2,$3,$4)", order.ID, item.ProductID, item.Quantity, item.UnitPrice)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to insert Item Order"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, SuccessResponse{Success: "Order created successfully"})
}

// ListOrdersUser godoc
// @Summary Lista os pedidos de um usuário
// @Description Retorna todos os pedidos associados ao usuário autenticado
// @Tags Pedidos
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse "Error to find Order"
// @Router /orders [get]
// @Security BearerAuth
func ListOrdersUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := config.DB.Query("SELECT id FROM pedidos WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		o.UserID = userID
		rows.Scan(&o.ID)
		items, _ := config.DB.Query("SELECT product_id, quantity, unit_price FROM itens_pedido WHERE order_id = $1", o.ID)
		for items.Next() {
			var item models.OrderItems
			items.Scan(&item.ProductID, &item.Quantity, &item.UnitPrice)
			o.Items = append(o.Items, item)
		}
		orders = append(orders, o)
	}
	c.JSON(http.StatusOK, orders)
}
