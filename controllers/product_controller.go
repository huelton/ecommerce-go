package controllers

import (
	"ecommerce/config"
	"ecommerce/models"
	"ecommerce/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var productRepo = repositories.NewProductRepository(config.DB)

// @Summary Crate a product
// @Description Create a product
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {object} models.Product
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
	err := productRepo.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func ListProducts(c *gin.Context) {
	products, err := productRepo.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to find products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context) {
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Just Admins should update products"})
		return
	}
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	// Alterando o parametro para int
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	product.ID = int(id)

	err = productRepo.UpdateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}
