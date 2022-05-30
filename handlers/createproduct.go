package handlers

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func (h *Handler) CreateProducts(c *gin.Context) {
	user, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusInternalServerError, "not logged in")
	}

	seller := user.(*models.Seller)

	product := models.Product{}

	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to bind product",
		})
		return
	}

	product.SellerId = seller.ID
	fmt.Println(product.SellerId)

	err = h.DB.CreateProduct(product)
	if err != nil {
		log.Println("check error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error creating product in handler"})
		return
	}

	response := gin.H{

		"data":      product,
		"status":    http.StatusText(http.StatusCreated),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusCreated, response)
}
