package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//GET ALL PRODUCTS FROM DB
func (h *Handler) GetAllProducts(c *gin.Context) {
	products := h.DB.GetAllProducts()

	c.JSON(http.StatusFound, gin.H{
		"All Products": products,
	})
}
