package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ListAllOrdersAdmin godoc
// @Summary Lista os pedidos de todos os usuarios Admin
// @Description Retorna todos os pedidos dos usuarios autenticado como admin
// @Tags Dashboard Admin
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} ErrorResponse "Error to find Order"
// @Router /admin/orders [get]
// @Security BearerAuth
func ListAllOrdersAdmin(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should access this route"})
		return
	}
	rows, err := config.DB.Query("SELECT id, status FROM pedidos ORDER BY create_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Error to find Order"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status); err != nil {
			continue
		}
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

func CountAllOrdersAdmin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	/*isAdmin := c.GetBool("is_admin")
	if !isAdmin {
	c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should access this route"})
	return
	}*/
	var count models.CountOrders
	var querie = "SELECT COUNT(CASE WHEN status = 'cancelado' THEN 1 END) AS cancelados, COUNT(CASE WHEN status = 'pendente' THEN 1 END) AS pendentes, COUNT(CASE WHEN status = 'pago' THEN 1 END) AS pagos, (SELECT SUM(ip.quantity * ip.unit_price) FROM itens_pedido ip JOIN pedidos p2 ON p2.id = ip.order_id WHERE p2.status = 'pago') AS total_valor_pago  FROM pedidos p"
	err := config.DB.QueryRow(querie).Scan(&count.Canceled, &count.Pending, &count.Paid, &count.TotalPaid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find count orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Count Orders",
		"count":   count,
	})
}
