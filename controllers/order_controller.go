package controllers

import (
	"database/sql"
	"ecommerce/config"
	"ecommerce/models"
	"fmt"
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
		var stockQuantity int
		err := tx.QueryRow("SELECT price, quantity FROM produtos WHERE id = $1", item.ProductID).Scan(&price, &stockQuantity)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Product"})
			return
		}

		if item.Quantity > stockQuantity {
			tx.Rollback()
			message := fmt.Sprintf("Stock is not enough for this product: %d", item.ProductID)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: message})
			return
		}

		_, err = tx.Exec("UPDATE produtos SET quantity = quantity - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to update Stock"})
			return
		}

		_, err = tx.Exec("INSERT INTO itens_pedido (order_id, product_id, quantity, unit_price) VALUES ($1,$2,$3,$4)", order.ID, item.ProductID, item.Quantity, price)
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
	rows, err := config.DB.Query("SELECT id, status FROM pedidos WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order"})
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

// OrderPayment realiza o pagamento de um pedido
// @Summary Realiza o pagamento de um pedido
// @Description Atualiza o status de um pedido para "pago", se ele ainda não estiver pago
// @Tags Pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do Pedido"
// @Success 200 {object} SuccessResponse "Pagamento realizado com sucesso"
// @Failure 400 {object} ErrorResponse "Pedido já está pago"
// @Failure 404 {object} ErrorResponse "Pedido não encontrado"
// @Failure 500 {object} ErrorResponse "Erro interno ao processar o pedido"
// @Router /orders/{id}/payment [put]
// @Security BearerAuth
func OrderPayment(c *gin.Context) {
	userID := c.GetInt("user_id")
	orderID := c.Param("id")

	var status string
	err := config.DB.QueryRow("SELECT status FROM pedidos WHERE id = $1 and user_id = $2", orderID, userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order"})
		return
	}

	if status == "pago" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Order Already paid"})
		return
	}

	_, err = config.DB.Exec("UPDATE pedidos SET status 'pago' WHERE id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to update Status order"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Success: "Payment Order created successfully"})
}

// CancelOrder realiza o cancelamento de um pedido
// @Summary Realiza o cancelamento de um pedido
// @Description Atualiza o status de um pedido para "cancelado", se ele ainda não estiver cancelado
// @Tags Pedidos
// @Accept json
// @Produce json
// @Param id path string true "ID do Pedido"
// @Success 200 {object} SuccessResponse "Cancelamento realizado com sucesso"
// @Failure 400 {object} ErrorResponse "Pedido já está pago ou cancelado"
// @Failure 404 {object} ErrorResponse "Pedido não encontrado"
// @Failure 500 {object} ErrorResponse "Erro interno ao cancelar o pedido"
// @Router /orders/{id}/cancel [put]
// @Security BearerAuth
func CancelOrder(c *gin.Context) {
	userID := c.GetInt("user_id")
	orderID := c.Param("id")

	var status string
	err := config.DB.QueryRow("SELECT status FROM pedidos WHERE id = $1 and user_id = $2", orderID, userID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order"})
		}
		return
	}

	if status != "pendente" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Order Already paid or cancel"})
		return
	}

	//get items to restore stock
	rows, err := config.DB.Query("SELECT product_id, quantity FROM itens_pedido WHERE order_id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order Itens"})
		return
	}
	defer rows.Close()

	type item struct {
		ProductID int
		Quantity  int
	}

	var itens []item
	for rows.Next() {
		var i item
		rows.Scan(&i.ProductID, &i.Quantity)
		itens = append(itens, i)
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to begin the Transaction"})
		return
	}

	for _, i := range itens {
		_, err = tx.Exec("UPDATE produtos SET quantity = quantity + $1 WHERE id = $2", i.Quantity, i.ProductID)
		_, err = tx.Exec("UPDATE itens_pedido SET quantity = 0, unit_price = 0 WHERE product_id = $1 AND order_id = $2", i.ProductID, orderID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to restore Stock"})
			return
		}
	}

	_, err = tx.Exec("UPDATE pedidos SET status ='cancelado' WHERE id = $1", orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to cancel order"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, SuccessResponse{Success: "Order cancel with successfully"})
}
