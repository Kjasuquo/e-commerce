package handlers

import (
	"github.com/decadevs/shoparena/database"
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Handler struct {
	DB       database.DB
	Mail     database.Mailer
	Paystack database.Paystack
}

func PingHandler(c *gin.Context) {
	// healthcheck
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) SellerUpdatePassword(c *gin.Context) {
	var password models.PasswordResetReq
	userI, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{"message": "internal server get error"})
	}
	user, ok := userI.(*models.Seller)
	if !ok {
		c.JSON(500, gin.H{"message": "invalid assert"})
	}
	seller, err := h.DB.FindSellerByEmail(user.Email)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		c.Abort()
		return
	}
	err = c.BindJSON(&password)
	if err != nil {
		c.JSON(406, gin.H{"message": "Provide all Relevant Fields"})
		c.Abort()
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(seller.PasswordHash), []byte(password.OldPassword))
	if err != nil {
		c.JSON(500, gin.H{"message": "OldPassword is incorrect"})
		c.Abort()
		return
	}

	if password.NewPassword != password.ConfirmNewPassword {
		c.JSON(400, gin.H{"message": "Password Mismatch"})
		c.Abort()
		return
	}
	passwordhash, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	_, err = h.DB.SellerUpdatePassword(seller.PasswordHash, string(passwordhash))
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"message": "successfully reset password, please login"})
}

func (h *Handler) BuyerUpdatePassword(c *gin.Context) {
	var password models.PasswordResetReq
	userI, exists := c.Get("user")
	log.Println("testtttt", userI)
	if !exists {
		c.JSON(500, gin.H{"message": "internal server get error"})
	}
	user, ok := userI.(*models.Buyer)
	if !ok {
		c.JSON(500, gin.H{"message": "invalid assert"})
	}

	err := c.BindJSON(&password)
	if err != nil {
		c.JSON(406, gin.H{"message": "Provide all Relevant Fields"})
		c.Abort()
		return
	}

	if password.NewPassword != password.ConfirmNewPassword {
		c.JSON(400, gin.H{"message": "Password Mismatch"})
		c.Abort()
		return
	}

	buyer, err := h.DB.FindBuyerByEmail(user.Email)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(buyer.PasswordHash), []byte(password.OldPassword))
	if err != nil {
		c.JSON(500, gin.H{"message": "OldPassword is incorrect"})
		c.Abort()
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = h.DB.BuyerUpdatePassword(buyer.PasswordHash, string(passwordHash))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message": "successfully reset password, please login"})
}

/*
func (h *Handler) SellerUpdatePasswordHandler(c *gin.Context) {
	var password models.PasswordResetReq
	email := c.Param("email")
	seller, err := h.DB.FindSellerByEmail(email)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		c.Abort()
		return
	}
	err = c.BindJSON(&password)
	if err != nil {
		c.JSON(406, gin.H{"message": "Provide all Relevant Fields"})
		c.Abort()
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(seller.PasswordHash), []byte(password.OldPassword))
	if err != nil {
		c.JSON(500, gin.H{"message": "OldPassword is incorrect"})
		c.Abort()
		return
	}

	if password.NewPassword != password.ConfirmNewPassword {
		c.JSON(400, gin.H{"message": "Password Mismatch, Please make sure newpassword & confirmNewPassword match"})
		c.Abort()
		return
	}
	passwordhash, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	_, err = h.DB.SellerUpdatePassword(seller.PasswordHash, string(passwordhash))
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(200, gin.H{"message": "successfully reset password"})
}

func (h *Handler) BuyerUpdatePasswordHandler(c *gin.Context) {
	var password models.PasswordResetReq

	email := c.Param("email")

	err := c.BindJSON(&password)
	if err != nil {
		c.JSON(406, gin.H{"message": "Provide all Relevant Fields"})
		c.Abort()
		return
	}

	if password.NewPassword != password.ConfirmNewPassword {
		c.JSON(400, gin.H{"message": "Password Mismatch"})
		c.Abort()
		return
	}

	buyer, err := h.DB.FindBuyerByEmail(email)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		c.Abort()
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(buyer.PasswordHash), []byte(password.OldPassword))
	if err != nil {
		c.JSON(500, gin.H{"message": "OldPassword is incorrect"})
		c.Abort()
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = h.DB.BuyerUpdatePassword(buyer.PasswordHash, string(passwordHash))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"message": "internal server error"})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"message": "successfully reset password"})
}

*/
