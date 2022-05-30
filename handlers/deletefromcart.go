package handlers

import (
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteFromCart(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, "error getting valid user from token")
		return
	}
	buyer := user.(*models.Buyer)
	StringID := c.Param("id")
	cartProductID, _ := strconv.Atoi(StringID)

	err := h.DB.DeleteCartProduct(buyer.ID, uint(cartProductID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error deleting from DB")
		return
	}

	c.JSON(http.StatusOK, "successfully deleted")

}

func (h *Handler) DeleteAllCartProducts(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, "error getting valid user from token")
		return
	}

	buyer := user.(*models.Buyer)
	err := h.DB.DeleteAllFromCart(buyer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error deleting all products from cart")
		return
	}

	c.JSON(http.StatusOK, "successfully deleted")
}
