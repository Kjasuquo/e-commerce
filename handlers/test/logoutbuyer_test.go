package test

import (
	"fmt"
	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestBuyerLogout(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(err.Error())
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{DB: mockDB}
	route, _ := router.SetupRouter(h)
	accClaim, _ := services.GenerateClaims("mike123@gmail.com")
	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	//ref, err := services.GenerateToken(jwt.SigningMethodHS256, refClaim, &secret)
	if err != nil {
		t.Fail()
	}
	buyer := &models.Buyer{}
	buyer.Username = "mike123"
	buyer.Email = "mike123@gmail.com"
	buyer.Token = *acc

	mockDB.EXPECT().TokenInBlacklist(acc).Return(false).Times(1)
	mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(buyer, nil).Times(1)
	mockDB.EXPECT().AddTokenToBlacklist(buyer.Email, buyer.Token).Return(nil).Times(1)

	ht := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/buyer/logout", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
	route.ServeHTTP(ht, req)
	assert.Equal(t, http.StatusOK, ht.Code)
	assert.Contains(t, ht.Body.String(), "successfully")

	t.Run("test bad request", func(t *testing.T) {
		buyer := &models.Seller{}
		buyer.Username = "mike123"
		buyer.Email = "mike123@gmail.com"
		buyer.Token = *acc

		ht := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/buyer/logout", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "hdhhdhddhdh"))
		route.ServeHTTP(ht, req)
		assert.Equal(t, http.StatusUnauthorized, ht.Code)
		//assert.Contains(t, ht.Body.String(), "authorize access token error")
	})
}
