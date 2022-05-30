package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(dbPass, pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(pass)) == nil
}

func HashPassword(pass string) (string, error) {
	bytePass := []byte(pass)
	hPass, err := bcrypt.GenerateFromPassword(bytePass, 12)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	hashpass := string(hPass)
	return hashpass, nil
}

func (h *Handler) LoginBuyerHandler(c *gin.Context) {
	buyer := &models.Buyer{}
	buyerLoginRequest := &struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&buyerLoginRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "bad request"})
		return
	}

	hashpass, err := HashPassword(buyerLoginRequest.Password)
	if err != nil {
		//write error
		c.JSON(500, gin.H{"message": "internal Server Error"})
		return
	}
	buyer, sqlErr := h.DB.FindBuyerByEmail(buyerLoginRequest.Email)

	if sqlErr != nil {
		fmt.Println(sqlErr)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	}
	fmt.Println(hashpass)

	if err := bcrypt.CompareHashAndPassword([]byte(buyer.PasswordHash), []byte(buyerLoginRequest.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid Password"})
		return
	}

	// Generates access claims and refresh claims
	accessClaims, refreshClaims := services.GenerateClaims(buyer.Email)

	secret := os.Getenv("JWT_SECRET")
	accToken, err := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		log.Printf("token generation error err: %v\n", err)
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}

	refreshToken, err := services.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		log.Printf("token generation error err: %v\n", err)
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}
	c.Header("refresh_token", *refreshToken)
	c.Header("access_token", *accToken)
	response.JSON(c, "login successful", http.StatusOK, gin.H{
		"user":          buyer,
		"access_token":  *accToken,
		"refresh_token": *refreshToken,
	}, nil)

}
