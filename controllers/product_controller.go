package controllers

import (
	"net/http"

	"ecommerce/config"
	"ecommerce/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a product
// @Description Return a product created
// @Security BearerAuth
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Dados do produto"
// @Success 201 {object} models.Product
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should create products"})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}
	err := config.DB.QueryRow("INSERT INTO produtos(name, description, price, quantity) VALUES ($1,$2,$3,$4) RETURNING id",
		product.Name, product.Description, product.Price, product.Quantity,
	).Scan(&product.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)

}

func ListProducts(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, description, price, quantity FROM produtos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find products"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Quantity)
		if err != nil {
			continue
		}
		products = append(products, p)
	}
	c.JSON(http.StatusOK, products)
}
