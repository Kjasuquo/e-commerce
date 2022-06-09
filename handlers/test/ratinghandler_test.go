package test

import (
	"encoding/json"
	"fmt"
	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSelerRating(t *testing.T) {
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

	rating := struct {
		Rating uint `json:"rating"`
		Id     uint `json:"Id"`
	}{
		5,
		1,
	}

	buyer := models.Buyer{
		User: models.User{
			Email: "chuks@gmail.com",
		},
		Orders: nil,
	}
	seller := models.Seller{
		User: models.User{
			Email: "jdoe@gmail.com",
		},
		Orders: nil,
	}
	updateData := models.UpdateRating{
		5,
		5,
		1,
	}
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindBuyerByEmail("chuks@gmail.com").Return(&buyer, nil)
	mockDB.EXPECT().FindSellerById(uint(1)).Return(&seller, nil)
	mockDB.EXPECT().UpdateSellerRating(uint(0), &updateData).Return(nil)

	ratingPayload, err := json.Marshal(rating)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/buyer/rateaseller",
		strings.NewReader(string(ratingPayload)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
	route.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "feedback", "thankyou")
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestProductRating(t *testing.T) {
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

	rating := struct {
		Rating uint `json:"rating"`
		Id     uint `json:"Id"`
	}{
		5,
		1,
	}

	buyer := models.Buyer{
		User: models.User{
			Email: "chuks@gmail.com",
		},
		Orders: nil,
	}
	product := models.Product{
		Model:                   gorm.Model{},
		SellerId:                1,
		CategoryId:              1,
		Category:                models.Category{},
		Title:                   "shoe",
		Description:             "running shoes",
		Price:                   7000,
		Images:                  nil,
		Rating:                  5,
		TotalRatings:            5,
		NumberOfRatingsReceived: 1,
		Quantity:                10,
	}
	updateData := models.UpdateRating{
		5,
		10,
		2,
	}
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindBuyerByEmail("chuks@gmail.com").Return(&buyer, nil)
	mockDB.EXPECT().FindProductById(uint(1)).Return(&product, nil)
	mockDB.EXPECT().UpdateProductRating(uint(0), &updateData).Return(nil)

	ratingPayload, err := json.Marshal(rating)
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/buyer/rateaproduct",
		strings.NewReader(string(ratingPayload)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
	route.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "feedback", "thankyou")
	assert.Equal(t, w.Code, http.StatusOK)
}
