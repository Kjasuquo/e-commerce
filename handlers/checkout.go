package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Fees struct {
	Amount float64 `json:"amount" binding:"required"`
}

type Transaction struct {
	UserID      uint    `json:"user_id"`
	Amount      float64 `json:"amount"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	CallBackUrl string  `json:"callback_url"`
	Reference   string  `json:"reference"`
}

func (h *Handler) Pay(c *gin.Context) {
	userI, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "please log in"})
		return
	}
	user := userI.(*models.Buyer)

	var Incoming Fees

	err := c.BindJSON(&Incoming)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	secret := os.Getenv("JWT_SECRET")

	newClaims := jwt.MapClaims{
		"email":   user.Email,
		"cart_id": user.ID,
		"u_id":    uuid.NewString(),
	}

	token, err := services.GenerateToken(jwt.SigningMethodHS256, newClaims, &secret)
	if err != nil {
		log.Printf("new token generation error err: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	transaction := Transaction{
		UserID:      user.ID,
		Amount:      Incoming.Amount * 100,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		CallBackUrl: "https://oja-ecommerce.herokuapp.com/api/v1/callback",
		Reference:   *token,
	}
	m, err := json.Marshal(transaction)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "not marshalling"})
		return
	}
	authorizationUrl, err := h.Paystack.InitializePayment(m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "not valid"})
		return
	}

	response.JSON(c, "Transaction initialized", http.StatusOK, gin.H{
		"authorization_url": authorizationUrl,
	}, nil)

}

func (h *Handler) Callback(c *gin.Context) {
	reference := c.Query("reference")

	_, err := h.Paystack.VerifyReference(reference)
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, "https://shoparena-frontend-phi.vercel.app/buyer/payment/unsuccessful")
		return
	}

	//fmt.Printf("this is before the decode %v \n", resp)
	//claims, err := DecodeTokenForPayment(reference)

	secret := os.Getenv("JWT_SECRET")
	claims, err := h.Paystack.PayStackDecodeToken(reference, secret)
	if err != nil {
		c.Redirect(http.StatusFound, "https://shoparena-frontend-phi.vercel.app/buyer/payment/unsuccessful")
		return
	}

	var ID string
	for i, v := range claims {
		if i == "cart_id" {
			ID = fmt.Sprintf("%v", v)
		}
	}

	cartID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Printf("error converting stringed cartID to int")
		return
	}

	err = h.DB.DeletePaidFromCart(uint(cartID))
	if err != nil {
		log.Println(err)
		c.Redirect(http.StatusFound, "https://shoparena-frontend-phi.vercel.app/buyer/payment/unsuccessful")
		return
	}

	c.Redirect(http.StatusFound, "https://shoparena-frontend-phi.vercel.app/buyer/payment/successful")
	return
}
