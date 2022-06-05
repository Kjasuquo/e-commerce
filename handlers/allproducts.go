package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//GET ALL PRODUCTS FROM DB
func (h *Handler) GetAllProducts(c *gin.Context) {
	products := h.DB.GetAllProducts()

	c.JSON(http.StatusOK, gin.H{
		"All_Products": products,
	})
	//c.JSON(http.StatusOK, gin.H{
	//	"all_Products": products,
	//})
	c.JSON(http.StatusOK, products)
}
