package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

//GET TOTAL PRODUCT COUNT FOR A SELLER
func (h *Handler) GetTotalProductCountForSeller(c *gin.Context) {

	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	sellerID := seller.ID

	//find seller with the retrieved ID and return the seller and its preloaded product
	Seller, err := h.DB.FindIndividualSellerShop(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist ; Seller with this ID not found in product table",
			"error":   err.Error(),
		})
		return
	}

	//sum of the quantity of each product
	var total_quantity uint

	product := Seller.Product
	//total count based on the individual product
	for i := 0; i < len(product); i++ {
		total_quantity += product[i].Quantity
	}

	//total count based on an individual product
	c.IndentedJSON(http.StatusOK, gin.H{
		"Message":                fmt.Sprintf("Seller has %d different categories of product to sell", len(product)),
		"Seller_ID":              sellerID,
		"Product_Count":          len(product),
		"Product_Quantity":       total_quantity,
		"First_Product_Quantity": product[0].Quantity,
	})

}

func (h *Handler) SellerIndividualProduct(c *gin.Context) {
	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	sellerID := seller.ID

	product, err := h.DB.FindSellerProduct(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Message": "Error Exist Count not find individual seller product",
			"error":   err.Error(),
		})
		return
	}

	if len(product) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "No Product Found In Shop",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": "Seller Has Product In Shop",
		"Product": product,
	})

}
