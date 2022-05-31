package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) DeleteAllSellerProducts(c *gin.Context) {
	Seller, err := h.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	err = h.DB.DeleteAllSellerProducts(Seller.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting all seller products"})
		return
	}
	c.JSON(http.StatusOK, "all products deleted successfully")
}
