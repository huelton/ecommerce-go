package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
