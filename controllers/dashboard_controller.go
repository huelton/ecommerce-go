package controllers

import (
	"ecommerce/config"
	"ecommerce/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

var rep = repositories.NewOrderRepository(config.DB)

func CountAllOrdersAdmin(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	/*isAdmin := c.GetBool("is_admin")
	if !isAdmin {
	c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should access this route"})
	return
	}*/
	count, err := rep.GetOrderCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find count orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "Count Orders",
		"sum_total": count.TotalPaid,
		"count":     count,
	})
}
