package handlers

import (
	"net/http"

	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) BuyerSignUpHandler(c *gin.Context) {
	buyer := &models.Buyer{}
	err := c.ShouldBindJSON(buyer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to bind JSON",
		})
		return
	}
	if buyer.Username == "" || buyer.FirstName == "" || buyer.LastName == "" || buyer.Password == "" || buyer.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enter all fields",
		})
		return
	}
	validEmail := buyer.ValidMailAddress()
	if validEmail == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "enter valid email",
		})
		return
	}

	_, err = h.DB.FindBuyerByUsername(buyer.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username exists",
		})
		return
	}
	_, err = h.DB.FindBuyerByEmail(buyer.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email exists",
		})
		return
	}

	_, err = h.DB.FindBuyerByPhone(buyer.PhoneNumber)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number exists",
		})
		return
	}

	if err = buyer.HashPassword(); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	_, err = h.DB.CreateBuyer(buyer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create buyer",
		})
		return
	}
	cart := &models.Cart{BuyerID: buyer.ID}
	_, err = h.DB.CreateBuyerCart(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": error.Error(err),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sign Up Successful",
	})

}

func (h *Handler) SellerSignUpHandler(c *gin.Context) {

	seller := &models.Seller{}
	err := c.ShouldBindJSON(seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to bind json",
		})
		return
	}

	if seller.Username == "" || seller.FirstName == "" || seller.LastName == "" || seller.Password == "" || seller.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enter all fields",
		})
		return
	}
	validEmail := seller.ValidMailAddress()
	if validEmail == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "enter valid email",
		})
		return
	}

	_, err = h.DB.FindSellerByUsername(seller.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username exists",
		})
		return
	}
	_, err = h.DB.FindSellerByEmail(seller.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email exists",
		})
		return
	}

	_, err = h.DB.FindSellerByPhone(seller.PhoneNumber)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "phone number exists",
		})
		return

	}
	if err := seller.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	_, err = h.DB.CreateSeller(seller)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create seller",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sign Up Successful",
	})
}
