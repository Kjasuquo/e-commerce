package handlers

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//UPDATE PRODUCT BY ID
func (h *Handler) UpdateProduct(c *gin.Context) {
	user, err := h.GetUserFromContext(c)
	log.Println("user in context", user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	prod := c.Param("id")
	prodID, _ := strconv.Atoi(prod)
	prodIdUint := uint(prodID)

	var product = models.Product{}

	if err := c.BindJSON(&product); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error binding in updateProduct handler"})
		return
	}
	productInDb, errP := h.DB.GetProductByID(prodIdUint)
	fmt.Println(productInDb)
	if errP != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Product with ID does not exist"})
		return
	}

	fmt.Println(product.Title, product.Price)

	err = h.DB.UpdateProductByID(prodIdUint, product)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error updating product by ID"})
		return
	}
	response := gin.H{
		"message":         "product updated successfully",
		"status code":     http.StatusOK,
		"updated product": product,
	}

	c.JSON(http.StatusOK, response)
}
