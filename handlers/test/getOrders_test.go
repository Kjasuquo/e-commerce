package test

import (
	"encoding/json"
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
	_ "math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAllSellerOrders(t *testing.T) {
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
	////Declaring FAKE testing variables
	Id := uint(5)
	categoryID := uint(gofakeit.Number(1, 9))
	sellerID := uint(5)

	buyerID := uint(gofakeit.Number(1, 10))
	productID := uint(gofakeit.Number(1, 10))
	sellerFirstName := gofakeit.FirstName()
	sellerLastName := gofakeit.LastName()
	sellerUserName := gofakeit.Username()
	sellerPhone := gofakeit.Phone()
	productCategory := gofakeit.CarModel()
	productTitle := gofakeit.CarType()
	productDescriptn := gofakeit.CarType()
	productPrice := gofakeit.Price(1200, 1000000)
	convPrice := uint(productPrice)
	quantity := uint(gofakeit.Number(1, 100))

	address := "lagos"

	//instantiating the gorm model object/struct
	testGormModel := gorm.Model{
		ID:        Id,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	//instantiating the user model object/struct
	testUserModel := models.User{
		Model:       testGormModel,
		FirstName:   sellerFirstName,
		LastName:    sellerLastName,
		Email:       sellerEmail,
		Username:    sellerUserName,
		Address:     address,
		PhoneNumber: sellerPhone,
	}

	//instantiating the seller object/struct
	testSellerModel := models.Seller{
		Model:   testGormModel,
		User:    testUserModel,
		Product: nil,
		Orders:  nil,
	}
	//instantiating the buyer object/struct
	testBuyerModel := models.Buyer{
		Model:  testGormModel,
		User:   testUserModel,
		Orders: nil,
	}

	category := models.Category{
		testGormModel,
		productCategory,
	}
	product := models.Product{
		Model:       testGormModel,
		SellerId:    sellerID,
		CategoryId:  categoryID,
		Category:    category,
		Title:       productTitle,
		Description: productDescriptn,
		Price:       convPrice,
		Quantity:    quantity,
	}
	//instantiating the buyer object/struct

	orderOne := models.Order{
		testGormModel,
		sellerID,
		testSellerModel,
		buyerID,
		testBuyerModel,
		productID,
		product,
	}
	orderTwo := models.Order{
		testGormModel,
		sellerID,
		testSellerModel,
		buyerID,
		testBuyerModel,
		productID,
		product,
	}

	orderThree := models.Order{
		testGormModel,
		sellerID,
		testSellerModel,
		buyerID,
		testBuyerModel,
		productID,
		product,
	}

	testOrders := []models.Order{orderOne, orderTwo, orderThree}
	testUser := models.User{
		Model:        testGormModel,
		FirstName:    sellerFirstName,
		LastName:     sellerLastName,
		Email:        sellerEmail,
		Username:     sellerUserName,
		PasswordHash: sellerPhone,
	}
	productSlice := []models.Product{product}

	testSeller := models.Seller{
		Model:   testGormModel,
		User:    testUser,
		Product: productSlice,
	}

	bodyJSON, err := json.Marshal(testOrders)
	if err != nil {
		t.Fail()
	}

	//authentication and authorisation
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false)
	mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(&testSeller, nil)

	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockDB.EXPECT().GetAllSellerOrder(uint(5)).Return(testOrders, nil)

		mockDB.EXPECT().GetAllSellerOrder(uint(5)).Return(testOrders, nil).AnyTimes()

		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellerorders/", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)

	})

}
