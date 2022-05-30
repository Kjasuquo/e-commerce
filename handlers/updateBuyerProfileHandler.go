package handlers

import (

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetBuyerFromContext(c *gin.Context) (*models.Buyer, error) {
	userI, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("error getting user from context")
	}
	user, ok := userI.(*models.Buyer)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return user, nil
}

func (h *Handler) UpdateBuyerProfileHandler(c *gin.Context) {

	buyer, err := h.GetBuyerFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	var update models.UpdateUser
	email := buyer.Email

	if errs := h.Decode(c, &update); errs != nil {
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	buyer.Email = email
	buyer.UpdatedAt = time.Now()
	if err := h.DB.UpdateBuyerProfile(buyer.ID, &update); err != nil {
		log.Printf("update buyer error : %v\n", err)
		c.JSON(http.StatusInternalServerError, []string{"update buyer not successful"})
		return
	}
	c.JSON(http.StatusOK, []string{"buyer updated successfully!"})
}
