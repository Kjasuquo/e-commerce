package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) DeleteSellerProduct(c *gin.Context) {
	seller, err := h.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	productID := c.Param("id")
	productIdInt, err := strconv.Atoi(productID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(uint(productIdInt))

	err = h.DB.DeleteProduct(uint(productIdInt), seller.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error deleting product"})
		return
	}

	c.JSON(http.StatusOK, "successfully deleted")
}
