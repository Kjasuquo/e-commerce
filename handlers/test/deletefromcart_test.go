package test

import (
	"errors"
	"fmt"
	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestDeleteCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)
	mockPaystack := mock_database.NewMockPaystack(ctrl)

	h := &handlers.Handler{DB: mockDB, Mail: mockMail, Paystack: mockPaystack}

	route, _ := router.SetupRouter(h)

	buyer := models.Buyer{
		User: models.User{Email: "joseph@yahoo.com"},
	}
	buyer.ID = uint(1)

	secret := os.Getenv("JWT_SECRET")
	accessClaims, _ := services.GenerateClaims(buyer.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)

	t.Run("Testing for error in Delete from cart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().DeleteCartProduct(buyer.ID, uint(2)).Return(errors.New("error deleting from DB"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/deletefromcart/2", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error deleting from DB")
	})

	t.Run("Testing for Delete from cart working", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().DeleteCartProduct(buyer.ID, uint(2)).Return(nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/deletefromcart/2", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "successfully deleted")
	})

	t.Run("Testing for error in Delete all from cart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().DeleteAllFromCart(buyer.ID).Return(errors.New("error deleting all products from cart"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/deleteallcart", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error deleting all products from cart")
	})

	t.Run("Testing for Delete all from cart working", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().DeleteAllFromCart(buyer.ID).Return(nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/deleteallcart", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "successfully deleted")
	})
}
