package test

import (
	"encoding/json"
	"fmt"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBuyerUpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)
	h := &handlers.Handler{DB: mockDB, Mail: mockMail}

	route, _ := router.SetupRouter(h)
	accClaims, _ := services.GenerateClaims("chuks@gmail.com")
	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaims, &secret)
	if err != nil {
		t.Fail()
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	resetPassword := struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}{
		OldPassword:        "12345678",
		NewPassword:        "123456789",
		ConfirmNewPassword: "123456789",
	}

	buyer := models.Buyer{
		User: models.User{
			Email:        "chuks@gmail.com",
			PasswordHash: string(passwordHash),
		},
		Orders: nil,
	}
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindBuyerByEmail("chuks@gmail.com").Return(&buyer, nil).Times(2)
	mockDB.EXPECT().BuyerUpdatePassword(string(passwordHash), gomock.Any()).Return(&buyer, nil)

	resetPasswordPayload, err := json.Marshal(resetPassword)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/buyer/updatepassword",
		strings.NewReader(string(resetPasswordPayload)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
	route.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "reset", "password")
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestSellerUpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)
	h := &handlers.Handler{DB: mockDB, Mail: mockMail}

	route, _ := router.SetupRouter(h)
	accClaims, _ := services.GenerateClaims("chuks@gmail.com")
	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaims, &secret)
	if err != nil {
		t.Fail()
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	resetPassword := struct {
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password"`
	}{
		OldPassword:        "12345678",
		NewPassword:        "123456789",
		ConfirmNewPassword: "123456789",
	}

	seller := models.Seller{
		User: models.User{
			Email:        "chuks@gmail.com",
			PasswordHash: string(passwordHash),
		},
	}
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindSellerByEmail("chuks@gmail.com").Return(&seller, nil).Times(2)
	mockDB.EXPECT().SellerUpdatePassword(string(passwordHash), gomock.Any()).Return(&seller, nil)

	resetPasswordPayload, err := json.Marshal(resetPassword)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/seller/updatepassword",
		strings.NewReader(string(resetPasswordPayload)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
	route.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "reset", "password")
	assert.Equal(t, w.Code, http.StatusOK)
}
