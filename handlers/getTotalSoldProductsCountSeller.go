package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) GetTotalSoldProductCount(c *gin.Context) {

	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)
	sellerFirstName := seller.FirstName
	sellerLastName := seller.LastName

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	sellerID := seller.ID
	//FIND THE cart HANDLER FUNCTION AND CALL IT HERE
	cartProduct, err := h.DB.FindPaidProduct(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Message": "Error Exist is finding sold product",
			"error":   err.Error(),
		})
		return
	}

	//FIND THE ORDERS AND CALL THE FUNCTION HERE
	var soldProductCount uint

	for i := 0; i < len(cartProduct); i++ {
		soldProductCount += cartProduct[i].TotalQuantity
	}

	if soldProductCount == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message":      "No Product has been Purchased",
			"Product_sold": soldProductCount,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Message":      "Total Product Sold by " + sellerFirstName + " " + sellerLastName + " is",
		"Product_sold": soldProductCount,
	})

}
