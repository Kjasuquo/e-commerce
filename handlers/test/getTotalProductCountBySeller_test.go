package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/decadevs/shoparena/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetTotalProductCountForSeller(t *testing.T) {

	ctrl := gomock.NewController(t)

	//creates a new mock instance
	mockDB := mock_database.NewMockDB(ctrl)

	h := &handlers.Handler{
		DB: mockDB,
	}

	route, _ := router.SetupRouter(h)

	sellerEmail := gofakeit.Email()
	accClaim, _ := services.GenerateClaims(sellerEmail)

	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}

	////Declaring FAKE testing variables
	Id := uint(gofakeit.Number(1, 10))
	categoryID := uint(gofakeit.Number(1, 9))
	sellerID := uint(0)
	sellerFirstName := gofakeit.FirstName()
	sellerLastName := gofakeit.LastName()
	sellerUserName := gofakeit.Username()
	sellerPhone := gofakeit.Phone()
	sellerImage := gofakeit.ImageURL(200, 500)
	sellerStatus := gofakeit.Bool()
	productCategory := gofakeit.CarModel()
	productTitle := gofakeit.CarType()
	productDescriptn := gofakeit.CarType()
	productPrice := gofakeit.Price(1200, 1000000)
	convPrice := uint(productPrice)
	productImage := gofakeit.ImageURL(200, 500)
	rating := uint(rand.Intn(5))
	quantity := uint(gofakeit.Number(1, 100))
	token := ""

	//instantiating the gorm model object/struct
	testGormModel := gorm.Model{
		ID:        Id,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	category := models.Category{
		testGormModel,
		productCategory,
	}
	//image object
	image := models.Image{
		testGormModel,
		sellerID,
		productImage,
	}

	sliceOfImg := []models.Image{image}

	productOne := models.Product{
		testGormModel,
		sellerID,
		categoryID,
		category,
		productTitle,
		productDescriptn,
		convPrice,
		sliceOfImg,
		rating,
		quantity,
	}
	productTwo := models.Product{
		testGormModel,
		sellerID,
		categoryID,
		category,
		productTitle,
		productDescriptn,
		convPrice,
		sliceOfImg,
		rating,
		quantity,
	}

	productThree := models.Product{
		testGormModel,
		sellerID,
		categoryID,
		category,
		productTitle,
		productDescriptn,
		convPrice,
		sliceOfImg,
		rating,
		quantity,
	}
	testProducts := []models.Product{productOne, productTwo, productThree}
	//instantiating the seller model object/struct

	testUser := models.User{
		Model:        testGormModel,
		FirstName:    sellerFirstName,
		LastName:     sellerLastName,
		Email:        sellerEmail,
		Username:     sellerUserName,
		PasswordHash: sellerPhone,
		Image:        sellerImage,
		IsActive:     sellerStatus,
		Token:        token,
	}

	testSeller := models.Seller{
		User:    testUser,
		Product: testProducts,
		Rating:  int(rating),
	}

	bodyJSON, err := json.Marshal(testProducts)
	if err != nil {
		t.Fail()
	}

	//authentication and authorisation
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(&testSeller, nil).Times(2)

	t.Run("Testing for Bad/Wrong Request", func(t *testing.T) {
		mockDB.EXPECT().FindIndividualSellerShop(sellerID).Return(nil, errors.New("Error Exist"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/seller/total/product/count", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Error Exist")
	})

	t.Run("Testing for Successful Request", func(t *testing.T) {
		mockDB.EXPECT().FindIndividualSellerShop(sellerID).Return(&testSeller, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/seller/total/product/count", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), fmt.Sprintf("%d", len(testSeller.Product)))
	})

}
