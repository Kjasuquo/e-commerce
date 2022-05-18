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
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)
	mockPaystack := mock_database.NewMockPaystack(ctrl)

	h := &handlers.Handler{DB: mockDB, Mail: mockMail, Paystack: mockPaystack}

	route, _ := router.SetupRouter(h)

	category := models.Category{
		Model: gorm.Model{ID: 1},
		Name:  "shirt",
	}

	seller := models.Seller{
		Model: gorm.Model{ID: 2},
	}

	product := models.Product{
		Model:       gorm.Model{ID: 1},
		SellerId:    seller.ID,
		CategoryId:  category.ID,
		Title:       "big shirt",
		Description: "nice big shirt",
		Price:       50000,
		Images:      nil,
		Rating:      5,
		Quantity:    7,
	}

	prodJASON, _ := json.Marshal(product)

	buyer := models.Buyer{
		User: models.User{Email: "joseph@yahoo.com"},
	}
	buyer.ID = 3

	addedProducts := []models.CartProduct{
		{
			Model:         gorm.Model{ID: 1},
			CartID:        2,
			ProductID:     1,
			TotalPrice:    4000,
			TotalQuantity: 2,
			OrderStatus:   false,
			SellerId:      1,
		},

		{
			Model:         gorm.Model{ID: 1},
			CartID:        2,
			ProductID:     2,
			TotalPrice:    3000,
			TotalQuantity: 1,
			OrderStatus:   false,
			SellerId:      2,
		},
	}

	addedProductJASON, _ := json.Marshal(addedProducts)

	productDetails := []models.ProductDetails{
		{
			Name:     "shirt",
			Price:    addedProducts[0].TotalPrice,
			Quantity: addedProducts[0].TotalQuantity,
			Images:   nil,
		},
		{
			Name:     "trouser",
			Price:    addedProducts[1].TotalPrice,
			Quantity: addedProducts[1].TotalQuantity,
			Images:   nil,
		},
	}

	productDetailsJASON, _ := json.Marshal(productDetails)

	secret := os.Getenv("JWT_SECRET")
	accessClaims, _ := services.GenerateClaims(buyer.Email)
	accToken, _ := services.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)

	t.Run("Testing for error in AddToChart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().AddToCart(product, &buyer).Return(errors.New("error adding to cart"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/addtocart", strings.NewReader(string(prodJASON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "cannot add to cart")
	})

	t.Run("No error in AddToChart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)
		mockDB.EXPECT().AddToCart(product, &buyer).Return(nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/addtocart", strings.NewReader(string(prodJASON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "successfully added")
	})

	t.Run("error in GetProductFromCart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)

		mockDB.EXPECT().GetCartProducts(&buyer).Return(nil, errors.New("error getting products from cart"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/viewcart", strings.NewReader(string(addedProductJASON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error getting cart products")
	})

	t.Run("No error in GetProductFromCart", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindBuyerByEmail(buyer.Email).Return(&buyer, nil)

		mockDB.EXPECT().GetCartProducts(&buyer).Return(addedProducts, nil)
		mockDB.EXPECT().ViewCartProducts(addedProducts).Return(productDetails, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/viewcart", strings.NewReader(string(addedProductJASON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *accToken))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(productDetailsJASON))
	})
}
