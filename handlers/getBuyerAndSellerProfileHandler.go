package handlers

import (
	"log"
	"net/http"

	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/gin-gonic/gin"
)

// handleShowProfile returns user's details
func (h *Handler) GetBuyerProfileHandler(c *gin.Context) {

		if userI, exists := c.Get("user"); exists {
			if user, ok := userI.(*models.Buyer); ok {
				response.JSON(c, "buyer details retrieved correctly", http.StatusOK, user, nil)
				return
			}
		}
		log.Printf("can't get buyer from context\n")
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
	
}


func (h *Handler) GetSellerProfileHandler(c *gin.Context) {

	if userI, exists := c.Get("user"); exists {
		if user, ok := userI.(*models.Seller); ok {
			response.JSON(c, "seller details retrieved correctly", http.StatusOK, user, nil)
			return
		}
	}
	log.Printf("can't get seller from context\n")
	response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})

}