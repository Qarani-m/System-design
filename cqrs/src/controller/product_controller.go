package controller

import (
	"context"
	"main/model"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service *service.ProductService
}

// RegisterRoutes wires the routes for products
func (pc *ProductController) RegisterRoutes(r *gin.Engine) {
	// group them under /products if you want
	r.GET("/products", pc.GetProducts)
	r.POST("/products", pc.CreateProduct)
}

func NewProductController(s *service.ProductService) *ProductController {
	return &ProductController{service: s}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req model.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.service.Create(context.Background(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "product created", "product": req})
}

func (pc *ProductController) GetProducts(c *gin.Context) {
	products, err := pc.service.FindAll(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
