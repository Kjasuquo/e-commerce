package handlers

import (
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SellerAllProducts(c *gin.Context) {
	user, present := c.Get("user")
	if !present {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
	}
	seller := user.(*models.Seller)
	sellerProducts, err := h.DB.GetSellersProducts(seller.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get products",
		})
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"SellerProducts": sellerProducts,
	})
}
