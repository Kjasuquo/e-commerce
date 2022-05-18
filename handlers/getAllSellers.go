package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSellers(c *gin.Context) {

	Seller, err := h.DB.GetAllSellers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error Exist in Getting All Sellers",
			"error":   err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "All Seller Found",
		"Sellers": Seller,
	})

}
