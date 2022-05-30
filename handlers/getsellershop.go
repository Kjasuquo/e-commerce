package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) HandleGetSellerShopByProfileAndProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		//authorize and authenticate seller
		seller, err := h.GetUserFromContext(c)

		if err != nil {
			c.JSON(http.StatusInternalServerError, []string{"internal server error"})
			return
		}
		sellerID := seller.ID

		//find seller with the retrieved ID and return the seller and its product
		Seller, err := h.DB.FindIndividualSellerShop(sellerID)

		if err != nil {
			log.Println("Error finding information in database:", err)
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"Message": "Error Exist ; Seller with this ID not found in product table",
				"error":   err.Error(),
			})
			return
		}

		if Seller == nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"Message": "Seller Shop Not Found",
			})
			return
		}

		//5. return a json object of seller profile and product if found
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message":     "Found Seller Shop by Profile and Product",
			"Seller_Shop": Seller,
		})

	}
}
