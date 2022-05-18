package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
)

func TestGetSellerProfile(t *testing.T) {
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
	accClaim, _ := services.GenerateClaims("cecilia@yahoo.com")

	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}

	//Declaring the FAKE testing variables
	sellerFirstName := gofakeit.FirstName()
	sellerLastName := gofakeit.LastName()
	sellerEmail := "cecilia@yahoo.com"
	sellerUserName := gofakeit.Username()
	sellerPhone := gofakeit.Phone()
	sellerImage := gofakeit.ImageURL(200, 500)
	sellerStatus := gofakeit.Bool()

	productName := gofakeit.CarModel()
	productPrice := gofakeit.Price(2500, 5000)
	price := uint(productPrice)
	productCat := gofakeit.CarModel()
	sellerID := uint(gofakeit.Number(0, 10))
	//testShopID := strconv.Itoa(int(sellerID))
	token := ""
	productImage := gofakeit.ImageURL(200, 500)

	//instantiating the gorm model object/struct
	testGormModel := gorm.Model{
		ID:        sellerID,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	//product category object
	category := models.Category{
		testGormModel,
		productCat,
	}

	//image object
	image := models.Image{
		testGormModel,
		sellerID,
		productImage,
	}

	sliceOfImg := []models.Image{image}

	//instantiating the user model object/struct
	testUser := models.User{
		testGormModel,
		sellerFirstName,
		sellerLastName,
		sellerEmail,
		sellerUserName,
		"",
		"",
		sellerPhone,
		"",
		"",
		sellerImage,
		sellerStatus,
		token,
	}

	products := []models.Product{
		{testGormModel,
			sellerID,
			sellerID,
			category,
			productCat,
			productName,
			price,
			sliceOfImg,
			uint(rand.Intn(20)),
			sellerID,
		},
		{testGormModel,
			sellerID,
			sellerID,
			category,
			productCat,
			"",
			price,
			sliceOfImg,
			uint(rand.Intn(5)),
			sellerID,
		},
	}

	//instantiating the seller model object/struct
	testSeller := models.Seller{
		User:    testUser,
		Product: products,
		Rating:  5,
	}

	bodyJSON, err := json.Marshal(testSeller)
	if err != nil {
		t.Fail()
	}

	t.Run("Test for error", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(nil, errors.New("an error occured"))
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/getsellerprofile",
			strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})

	t.Run("Test for successful profile retrieval", func(t *testing.T) {
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
		mockDB.EXPECT().FindSellerByEmail(sellerEmail).Return(&testSeller, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/getsellerprofile",
			strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "seller details retrieved correctly")
	})
}