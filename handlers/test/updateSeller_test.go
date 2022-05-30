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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestUpdateSellerDetailsHandler(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(err.Error())
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{DB: mockDB}
	route, _ := router.SetupRouter(h)
	accClaim, _ := services.GenerateClaims("ceciliaorji@yahoo.com")
	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}
	updateSeller := &models.UpdateUser{
		FirstName:   "Gabby",
		Address:     "Edo Tech Park",
		PhoneNumber: "123456789",
		Email:       "a@a.com",
		LastName:    "Gabby",
	}
	seller := models.Seller{
		User: models.User{
			Model: gorm.Model{
				ID: 23,
			},
			FirstName:   "Orji",
			LastName:    "Cecilia",
			PhoneNumber: "09052755438",
			Email:       "ceciliaorji@yahoo.com",
			Address:     "Edo Tech Park",
		},
		Orders: nil,
	}
	updatedSeller, err := json.Marshal(updateSeller)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindSellerByEmail(seller.Email).Return(&seller, nil).Times(2)

	t.Run("Test for error", func(t *testing.T) {
		mockDB.EXPECT().UpdateSellerProfile(seller.ID, updateSeller).Return(errors.New("error exist"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/updatesellerprofile",
			strings.NewReader(string(updatedSeller)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "update seller not successful")
	})

	t.Run("Test for successful update", func(t *testing.T) {
		mockDB.EXPECT().UpdateSellerProfile(seller.ID, updateSeller).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/updatesellerprofile",
			strings.NewReader(string(updatedSeller)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		log.Println(w.Body.String())
		assert.Contains(t, w.Body.String(), "seller updated successfully!")
	})
}
