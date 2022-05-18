package test

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "shoparena234das")
	code := m.Run()
	os.Exit(code)
}

func TestUpdateBuyerDetailsHandler(t *testing.T) {
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

	updateBuyer := &models.UpdateUser{
		FirstName:   "Gabby",
		Address:     "Edo Tech Park",
		PhoneNumber: "123456789",
		Email:       "a@a.com",
		LastName:    "Gabby",
	}
	buyer := models.Buyer{
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
	updatedBuyer, err := json.Marshal(updateBuyer)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil).Times(2)

	t.Run("Test for error", func(t *testing.T) {
		mockDB.EXPECT().UpdateBuyerProfile(buyer.ID, updateBuyer).Return(errors.New("error exist"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/updatebuyerprofile",
			strings.NewReader(string(updatedBuyer)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		log.Println(w.Body.String())
		assert.Contains(t, w.Body.String(), "update buyer not successful")
	})

	t.Run("Test for successful update", func(t *testing.T) {
		mockDB.EXPECT().UpdateBuyerProfile(buyer.ID, updateBuyer).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/updatebuyerprofile",
			strings.NewReader(string(updatedBuyer)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "buyer updated successfully!")
	})

}
