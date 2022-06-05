package handlers

import (
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) HandleGetSellerShopByProfileAndProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		sID := c.Param("id")
		sellerID, _ := strconv.Atoi(sID)

		//find seller with the retrieved ID and return the seller and its product
		Seller, err := h.DB.FindIndividualSellerShop(uint(sellerID))

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
		var Shop []models.Seller
		Shop = append(Shop, *Seller)

		//5. return a json object of seller profile and product if found
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message":     "Found Seller Shop by Profile and Product",
			"Seller_Shop": Shop,
		})

	}
}
