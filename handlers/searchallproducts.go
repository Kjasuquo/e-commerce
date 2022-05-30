package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) SearchProductHandler(c *gin.Context) {
	categoryName := c.Query("category")
	lowerPrice := c.Query("lower-price")
	upperPrice := c.Query("upper-price")
	name := c.Query("name")

	product, err := h.DB.SearchProduct(lowerPrice, upperPrice, categoryName, name)
	if err != nil {
		log.Println("handler error in search product", err)
		c.JSON(http.StatusInternalServerError, "Product Not found")
		return
	}

	if len(product) == 0 {
		c.JSON(http.StatusInternalServerError, "no such product")
		return
	}

	c.JSON(http.StatusOK, product)
}
