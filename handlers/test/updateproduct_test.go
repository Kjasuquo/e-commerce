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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestUpdateProduct(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println(err.Error())
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{
		DB: mockDB,
	}
	route, _ := router.SetupRouter(h)

	seller := models.Seller{
		User: models.User{
			Model: gorm.Model{
				ID: 23,
			},
			FirstName:   "Victor",
			LastName:    "Ihemadu",
			PhoneNumber: "08160967596",
			Email:       "victorihemadu@gmail.com",
			Address:     "Edo Tech Park",
		},
		Orders: nil,
	}
	accClaim, _ := services.GenerateClaims(seller.Email)
	secret := os.Getenv("JWT_SECRET")

	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}

	category := models.Category{
		Model: gorm.Model{ID: 1},
		Name:  "Shirt",
	}

	product := &models.Product{
		Model:       gorm.Model{ID: 4, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
		SellerId:    2,
		CategoryId:  category.ID,
		Title:       "household",
		Description: "home items",
		Price:       78,
		Rating:      8,
		Quantity:    4903,
	}

	prod := models.Product{
		Model:       gorm.Model{ID: 4},
		SellerId:    2,
		CategoryId:  category.ID,
		Title:       "house",
		Description: "items",
		Price:       78,
		Rating:      8,
		Quantity:    4903,
	}
	prodJSON, err := json.Marshal(prod)
	if err != nil {
		t.Fail()
	}

	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindSellerByEmail(seller.Email).Return(&seller, nil)

	t.Run("Testing for valid update", func(t *testing.T) {
		mockDB.EXPECT().GetProductByID(product.ID).Return(product, nil)
		mockDB.EXPECT().UpdateProductByID(uint(4), prod)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPut,
			"/api/v1/update/product/4",
			strings.NewReader(string(prodJSON)))
		if err != nil {
			fmt.Printf("errrr here %v \n", err)
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "product updated successfully")
	})
}
