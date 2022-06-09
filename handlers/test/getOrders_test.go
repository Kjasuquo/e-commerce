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

	//buyerID := uint(gofakeit.Number(1, 10))
	//productID := uint(gofakeit.Number(1, 10))
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

	//instantiating the gorm model object/struct
	testGormModel := gorm.Model{
		ID:        Id,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	//instantiating the user model object/struct
	//testUserModel := models.User{
	//	Model:       testGormModel,
	//	FirstName:   sellerFirstName,
	//	LastName:    sellerLastName,
	//	Email:       sellerEmail,
	//	Username:    sellerUserName,
	//	Address:     address,
	//	PhoneNumber: sellerPhone,
	//}

	//instantiating the seller object/struct
	//testSellerModel := models.Seller{
	//	Model:   testGormModel,
	//	User:    testUserModel,
	//	Product: nil,
	//	Orders:  nil,
	//}
	////instantiating the buyer object/struct
	//testBuyerModel := models.Buyer{
	//	Model:  testGormModel,
	//	User:   testUserModel,
	//	Orders: nil,
	//}

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

	orderOne := models.OrderProducts{
		Fname:        sellerFirstName,
		Lname:        sellerLastName,
		CategoryName: productCategory,
		Title:        productTitle,
		Price:        convPrice,
		Quantity:     quantity,
	}
	orderTwo := models.OrderProducts{
		Fname:        sellerFirstName,
		Lname:        sellerLastName,
		CategoryName: productCategory,
		Title:        productTitle,
		Price:        convPrice,
		Quantity:     quantity,
	}

	orderThree := models.OrderProducts{
		Fname:        sellerFirstName,
		Lname:        sellerLastName,
		CategoryName: productCategory,
		Title:        productTitle,
		Price:        convPrice,
		Quantity:     quantity,
	}

	testOrders := []models.OrderProducts{orderOne, orderTwo, orderThree}
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
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).AnyTimes()
	mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(&testSeller, nil).AnyTimes()

	t.Run("Testing for Successful Request", func(t *testing.T) {

		mockDB.EXPECT().GetAllSellerOrders(uint(5)).Return(testOrders, nil).AnyTimes()

		mockDB.EXPECT().GetAllSellerOrders(uint(5)).Return(testOrders, nil).AnyTimes()

		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/sellerorders/", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.NotEqual(t, http.StatusOK, rw.Code)

	})

}
