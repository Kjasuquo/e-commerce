package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeSeller(findSellerByEmail func(string) (*models.Seller, error), tokenInBlacklist func(*string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		var seller *models.Seller
		var error error
		secret := os.Getenv("JWT_SECRET")
		accToken := services.GetTokenFromHeader(c)
		accessToken, accessClaims, err := services.AuthorizeToken(&accToken, &secret)
		if err != nil {
			log.Printf("authorize access token error: %s\n", err.Error())
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if tokenInBlacklist(&accessToken.Raw) || IsTokenExpired(accessClaims) {
			rt := &struct {
				RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
			}{}

			if err := c.ShouldBindJSON(rt); err != nil {
				log.Printf("no refresh token in request body: %v\n", err)
				RespondAndAbort(c, "", http.StatusBadRequest, nil, []string{"unauthorized"})
				return
			}

			if tokenInBlacklist(&rt.RefreshToken) {
				log.Printf("refresh token is blacklisted: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			_, rtClaims, err := services.AuthorizeToken(&rt.RefreshToken, &secret)
			if err != nil {
				log.Printf("authorize refresh token error: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			if IsTokenExpired(rtClaims) {
				log.Printf("refresh token is expired")
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			if sub, ok := rtClaims["sub"].(float64); ok && sub != 1 {
				log.Printf("invalid refresh token, the sub claim isn't correct")
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			//generate a new access token, and rest its exp time
			accessClaims["exp"] = time.Now().Add(services.AccessTokenValidity).Unix()
			newAccessToken, err := services.GenerateTokenWithMethod(jwt.SigningMethodHS256, accessClaims, &secret)
			if err != nil {
				log.Printf("can't generate new access token: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"can't generate new access token"})
				return
			}
			RespondAndAbort(c, "new access token generated", http.StatusOK, gin.H{"access_token": *newAccessToken}, []string{"access token is invalid"})
			return
		}

		if email, ok := accessClaims["user_email"].(string); ok {
			if seller, error = findSellerByEmail(email); error != nil {
				if inactiveErr, ok := err.(response.InActiveUserError); ok {
					RespondAndAbort(c, "", http.StatusBadRequest, nil, []string{inactiveErr.Error()})
					return
				}
				log.Printf("find user by email error: %v\n", err)
				RespondAndAbort(c, "", http.StatusNotFound, nil, []string{"user not found"})
				return
			}
		} else {
			log.Printf("user email is not string\n")
			RespondAndAbort(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		// set the user and token as context parameters.
		c.Set("user", seller)
		c.Set("access_token", accessToken.Raw)
		// calling next handler
		c.Next()
	}
}



func AuthorizeBuyer(findBuyerByEmail func(string) (*models.Buyer, error), tokenInBlacklist func(*string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		var buyer *models.Buyer
		var error error
		secret := os.Getenv("JWT_SECRET")
		accToken := services.GetTokenFromHeader(c)
		accessToken, accessClaims, err := services.AuthorizeToken(&accToken, &secret)
		if err != nil {
			log.Printf("authorize access token error: %s\n", err.Error())
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if tokenInBlacklist(&accessToken.Raw) || IsTokenExpired(accessClaims) {
			rt := &struct {
				RefreshToken string `json:"refresh_token,omitempty" binding:"required"`
			}{}

			if err := c.ShouldBindJSON(rt); err != nil {
				log.Printf("no refresh token in request body: %v\n", err)
				RespondAndAbort(c, "", http.StatusBadRequest, nil, []string{"unauthorized"})
				return
			}

			if tokenInBlacklist(&rt.RefreshToken) {
				log.Printf("refresh token is blacklisted: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			_, rtClaims, err := services.AuthorizeToken(&rt.RefreshToken, &secret)
			if err != nil {
				log.Printf("authorize refresh token error: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			if IsTokenExpired(rtClaims) {
				log.Printf("refresh token is expired")
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			if sub, ok := rtClaims["sub"].(float64); ok && sub != 1 {
				log.Printf("invalid refresh token, the sub claim isn't correct")
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"refresh token is invalid"})
				return
			}

			//generate a new access token, and rest its exp time
			accessClaims["exp"] = time.Now().Add(services.AccessTokenValidity).Unix()
			newAccessToken, err := services.GenerateTokenWithMethod(jwt.SigningMethodHS256, accessClaims, &secret)
			if err != nil {
				log.Printf("can't generate new access token: %v\n", err)
				RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"can't generate new access token"})
				return
			}
			RespondAndAbort(c, "new access token generated", http.StatusOK, gin.H{"access_token": *newAccessToken}, []string{"access token is invalid"})
			return
		}

		if email, ok := accessClaims["user_email"].(string); ok {
			if buyer, error = findBuyerByEmail(email); error != nil {
				if inactiveErr, ok := err.(response.InActiveUserError); ok {
					RespondAndAbort(c, "", http.StatusBadRequest, nil, []string{inactiveErr.Error()})
					return
				}
				log.Printf("find user by email error: %v\n", err)
				RespondAndAbort(c, "", http.StatusNotFound, nil, []string{"user not found"})
				return
			}
		} else {
			log.Printf("user email is not string\n")
			RespondAndAbort(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
			return
		}

		// set the user and token as context parameters.
		c.Set("user", buyer)
		c.Set("access_token", accessToken.Raw)
		// calling next handler
		c.Next()
	}
}

func RespondAndAbort(c *gin.Context, message string, status int, data interface{}, errs []string) {
	response.JSON(c, message, status, data, errs)
	c.Abort()
}

func IsTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := claims["exp"].(float64); ok {
		return float64(time.Now().Unix()) > exp
	}
	return true
}
