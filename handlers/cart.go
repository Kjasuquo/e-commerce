package handlers

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddToCart(c *gin.Context) {
	user1, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting valid user from token"})
		return
	}
	user := user1.(*models.Buyer)

	var product models.Product

	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	err = h.DB.AddToCart(product, user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot add to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully added"})
}

func (h *Handler) ViewCartProducts(c *gin.Context) {

	user1, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting valid user from token"})
		return
	}
	user := user1.(*models.Buyer)

	cartProducts, err := h.DB.GetCartProducts(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting cart products"})
		return
	}

	productDetails, err := h.DB.ViewCartProducts(cartProducts)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error getting product details"})
		return
	}

	c.JSON(http.StatusOK, productDetails)

}
