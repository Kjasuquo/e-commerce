package handlers

import (
	"github.com/decadevs/shoparena/models"
	"net/http"

	//"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) AllBuyerOrders(c *gin.Context) {
	user, present := c.Get("user")
	if !present {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{
			"message": "you are not logged in"})
	}

	buyer := user.(*models.Buyer)

	buyerWithOrder, err := h.DB.GetAllBuyerOrder(buyer.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get order(s)",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Buyer_Order": buyerWithOrder,
	})
}

func (h *Handler) AllSellerOrders(c *gin.Context) {
	user, present := c.Get("user")
	if !present {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
	}
	seller := user.(*models.Seller)
	sellerWithOrder, err := h.DB.GetAllSellerOrder(seller.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get order(s)",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Seller_Order": sellerWithOrder,
	})

}

func (h *Handler) SellerTotalOrders(c *gin.Context) {
	user, present := c.Get("user")
	if !present {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you are not logged in"})
	}
	seller := user.(*models.Seller)

	orderCount, err := h.DB.GetAllSellerOrderCount(seller.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"error": "unable to get total order(s) count",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"OrderCount": orderCount,
	})

}
