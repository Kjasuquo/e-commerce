package test

import (
	"encoding/json"
	"errors"
	"fmt"
	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCheckout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)
	mockPaystack := mock_database.NewMockPaystack(ctrl)

	h := &handlers.Handler{DB: mockDB, Mail: mockMail, Paystack: mockPaystack}

	route, _ := router.SetupRouter(h)
	fmt.Println(route)

	buyer := models.Buyer{}
	buyer.ID = 3
	buyer.Email = "joseph@yahoo.com"
	buyer.FirstName = "Joseph"
	buyer.LastName = "Asuquo"

	secret := os.Getenv("JWT_SECRET")
	accessClaims, _ := services.GenerateClaims(buyer.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)

	newClaims := jwt.MapClaims{
		"email":   buyer.Email,
		"cart_id": buyer.ID,
		"u_id":    uuid.NewString(),
	}

	token, _ := services.GenerateToken(jwt.SigningMethodHS256, newClaims, &secret)

	Incoming := handlers.Fees{Amount: 2000}

	transaction := handlers.Transaction{
		UserID:      buyer.ID,
		Amount:      Incoming.Amount * 100,
		FirstName:   buyer.FirstName,
		LastName:    buyer.LastName,
		Email:       buyer.Email,
		CallBackUrl: "http://localhost:8085/api/v1/callback",
		Reference:   *token,
	}

	transJASON, _ := json.Marshal(transaction)

	t.Run("Testing for error in Initializing", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockPaystack.EXPECT().InitializePayment(gomock.Any()).Return("", errors.New("error in Initializing Payment"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/pay", strings.NewReader(string(transJASON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "not valid")
	})

}
