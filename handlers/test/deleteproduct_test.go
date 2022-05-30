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
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDeleteProduct(t *testing.T) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Unable to load env")
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	h := &handlers.Handler{
		DB: mockDB,
	}
	route, _ := router.SetupRouter(h)

	seller := &models.Seller{
		User: models.User{
			Email: "a@gmail.com",
		},
	}
	accClaims, _ := services.GenerateClaims(seller.Email)
	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaims, &secret)
	if err != nil {
		log.Println("Unable to generate token")
	}
	testGorm := gorm.Model{
		ID:        1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	product := models.Product{
		Model:       testGorm,
		SellerId:    1,
		CategoryId:  1,
		Title:       "Material",
		Description: "Building",
		Price:       670,
		Rating:      10,
		Quantity:    80,
	}

	bodyJSON, err := json.Marshal(product)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindSellerByEmail(seller.Email).Return(seller, nil).Times(2)
	t.Run("testing for error deleting product", func(t *testing.T) {
		mockDB.EXPECT().DeleteProduct(product.ID, seller.ID).Return(errors.New("error fetching product"))
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/deleteproduct/"+strconv.Itoa(int(product.ID)),
			strings.NewReader(string(bodyJSON)))
		if err != nil {
			fmt.Printf("error occured")
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "error deleting product")
	})
	t.Run("testing for success deleting product", func(t *testing.T) {
		mockDB.EXPECT().DeleteProduct(product.ID, seller.ID).Return(nil)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, "/api/v1/deleteproduct/"+strconv.Itoa(int(product.ID)),
			strings.NewReader(string(bodyJSON)))
		if err != nil {
			fmt.Printf("error occured")
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "successfully deleted")
	})
}
