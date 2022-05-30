package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUserFromContext(c *gin.Context) (*models.Seller, error) {
	userI, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("error getting user from context")
	}
	user, ok := userI.(*models.Seller)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return user, nil
}

func (h *Handler) UpdateSellerProfileHandler(c *gin.Context) {

	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	var update models.UpdateUser
	email := seller.Email

	if errs := h.Decode(c, &update); errs != nil {
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	seller.Email = email
	seller.UpdatedAt = time.Now()
	if err := h.DB.UpdateSellerProfile(seller.ID, &update); err != nil {
		log.Printf("update seller error : %v\n", err)
		c.JSON(http.StatusInternalServerError, []string{"update seller not successful"})
		return
	}
	c.JSON(http.StatusOK, []string{"seller updated successfully!"})
}
