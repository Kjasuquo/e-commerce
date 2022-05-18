package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetProductById(c *gin.Context) {
	id := c.Param("id")
	idx, _ := strconv.Atoi(id)

	product, err := h.DB.GetProductByID(uint(idx))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error Exist in Getting All Sellers",
			"error":   err.Error()})
		log.Println("Error in getting product", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "product gotten by ID",
		"Product": product,
	})
}
