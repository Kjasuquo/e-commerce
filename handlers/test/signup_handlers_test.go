package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBuyerSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDb := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{DB: mockDb}
	r, _ := router.SetupRouter(h)

	user := models.User{
		Model:        gorm.Model{},
		FirstName:    "deji",
		LastName:     "garber",
		Email:        "garber@gmail.com",
		Username:     "deji",
		PasswordHash: "",
		Address:      "ikeja",
		PhoneNumber:  "08022334455",
		Password:     "password",
		Image:        "linkup",
		IsActive:     true,
	}
	cart := models.Cart{
		BuyerID: user.ID,
	}

	buyer := models.Buyer{
		User: user,
		Cart: cart,
	}

	mockDb.EXPECT().FindBuyerByUsername(user.Username).Return(&buyer, nil)
	newUser, err := json.Marshal(user)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/buyersignup", strings.NewReader(string(newUser)))
	r.ServeHTTP(w, req)

	mockDb.EXPECT().FindBuyerByUsername(user.Username).Return(&buyer, nil)

	t.Run("Check if user name exists", func(t *testing.T) {

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/buyersignup", strings.NewReader(string(newUser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username exists")

	})

}

func TestSellerSignUpHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDb := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{DB: mockDb}
	r, _ := router.SetupRouter(h)

	user := models.User{
		Model:        gorm.Model{},
		FirstName:    "deji",
		LastName:     "garber",
		Email:        "garber@gmail.com",
		Username:     "deji",
		PasswordHash: "",
		Address:      "ikeja",
		PhoneNumber:  "08022334455",
		Password:     "password",

		Image:    "linkup",
		IsActive: true,
	}

	seller := models.Seller{
		User:   user,
		Rating: 5,
	}
	mockDb.EXPECT().FindSellerByUsername(user.Username).Return(&seller, nil)

	newUser, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/sellersignup", strings.NewReader(string(newUser)))
	r.ServeHTTP(w, req)

	mockDb.EXPECT().FindSellerByUsername(user.Username).Return(&seller, nil)

	t.Run("Check if user name exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/sellersignup", strings.NewReader(string(newUser)))
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "username exists")
	})
}
