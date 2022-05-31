package handlers

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func (h *Handler) HandleLogoutBuyer(c *gin.Context) {
	//create token blacklist struct
	tokenI, exists := c.Get("access_token")
	if !exists {
		response.JSON(c, "error getting access token",
			http.StatusBadRequest, nil, []string{"error getting access token"})
	}
	buyer, exists := c.Get("user")
	if !exists {
		response.JSON(c, "error getting user from context",
			http.StatusBadRequest, nil, []string{"error getting user from context"})
	}
	buyerMod := buyer.(*models.Buyer)
	tokenStr := tokenI.(string)

	token, err := jwt.ParseWithClaims(tokenStr, &services.Claims{}, func(t *jwt.Token) (interface{}, error) {

		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return os.Getenv("JWT_SECRET"), nil
	})
	if claims, ok := token.Claims.(*services.Claims); !ok && !token.Valid {
		response.JSON(c, "error inserting claims",
			http.StatusBadRequest, nil, []string{"Claims not valid type"})

	} else {
		claims.StandardClaims.ExpiresAt = time.Now().Add(-time.Hour).Unix()
	}

	err = h.DB.AddTokenToBlacklist(buyerMod.Email, tokenStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error inserting Token into database", "Error": err})

	}

	c.JSON(http.StatusOK, gin.H{"message": "signed out successfully"})
}
