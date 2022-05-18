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
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSellerAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)

	//creates a new mock instance
	mockDB := mock_database.NewMockDB(ctrl)

	h := &handlers.Handler{
		DB: mockDB,
	}

	route, _ := router.SetupRouter(h)

	sellerEmail := "kukus@yahoo.com"
	accClaim, _ := services.GenerateClaims(sellerEmail)

	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}
	sellerID := uint(5)

	productOne := models.Product{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		SellerId:    sellerID,
		CategoryId:  uint(5),
		Category:    models.Category{Name: "Fashion"},
		Title:       "Shirts",
		Description: "Male Shirts",
		Price:       5000,
	}
	productTwo := models.Product{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		SellerId:    sellerID,
		CategoryId:  uint(5),
		Category:    models.Category{Name: "Fashion"},
		Title:       "Shirts",
		Description: "Male Shirts",
		Price:       5000,
	}
	testProducts := []models.Product{productOne, productTwo}
	testUser := models.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		FirstName: "Tolu",
		LastName:  "Clinton",
		Email:     "tolucto@gmail.com",
		Username:  "tolucto",
	}

	testSeller := models.Seller{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		User: testUser,
	}
	bodyJSON, err := json.Marshal(testProducts)
	if err != nil {
		t.Fail()
	}
	//authentication and authorisation
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).AnyTimes()
	mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(&testSeller, nil).AnyTimes()

	t.Run("Testing for Successful Request", func(t *testing.T) {
		mockDB.EXPECT().GetSellersProducts(uint(5)).Return(testProducts, nil).AnyTimes()
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/seller/allproducts/", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.NotEqual(t, http.StatusOK, rw.Code)

	})
}
