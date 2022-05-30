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
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetRemainingProductCountSeller(t *testing.T) {

	ctrl := gomock.NewController(t)

	//creates a new mock instance
	mockDB := mock_database.NewMockDB(ctrl)

	h := &handlers.Handler{DB: mockDB}

	route, _ := router.SetupRouter(h)

	sellerEmail := gofakeit.Email()
	accClaim, _ := services.GenerateClaims(sellerEmail)

	secret := os.Getenv("JWT_SECRET")
	acc, err := services.GenerateToken(jwt.SigningMethodHS256, accClaim, &secret)
	if err != nil {
		t.Fail()
	}

	//Declaring FAKE testing variables
	cartID := uint(gofakeit.Number(1, 10))
	sellerID := uint(0)
	productID := uint(1)
	Id := uint(gofakeit.Number(1, 10))
	sellerFirstName := gofakeit.FirstName()
	sellerLastName := gofakeit.LastName()
	sellerUserName := gofakeit.Username()
	sellerPhone := gofakeit.Phone()
	sellerImage := gofakeit.ImageURL(200, 500)
	sellerStatus := gofakeit.Bool()

	testShopID := uint(gofakeit.Number(0, 0))
	totalQuantity := uint(gofakeit.Number(1, 10))
	orderStatus := gofakeit.Bool()
	productPrice := gofakeit.Price(1200, 1000000)
	productImage := gofakeit.ImageURL(200, 500)
	totalPrice := uint(productPrice)
	rating := uint(rand.Intn(5))
	token := ""
	productName := gofakeit.CarModel()
	price := uint(productPrice)
	productCat := gofakeit.CarModel()

	//instantiating the gorm model object/struct
	testGormModel := gorm.Model{
		ID:        Id,
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

	indProduct := models.Product{
		testGormModel,
		sellerID,
		sellerID,
		category,
		productCat,
		productName,
		price,
		sliceOfImg,
		uint(rand.Intn(20)),
		sellerID,
	}

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

	//instantiating the seller model object/struct
	testSeller := models.Seller{
		User:    testUser,
		Product: products,
		Rating:  int(rating),
	}

	//instantiating the seller model object/struct
	_ = models.CartProduct{
		testGormModel,
		cartID,
		1,
		totalPrice,
		totalQuantity,
		orderStatus,
		0,
	}

	testCartProductTwo := models.CartProduct{
		testGormModel,
		cartID,
		productID,
		totalPrice,
		totalQuantity,
		orderStatus,
		sellerID,
	}

	bodyJSON, err := json.Marshal(products)
	if err != nil {
		t.Fail()
	}

	//authentication and authorisation
	mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
	mockDB.EXPECT().FindSellerByEmail(testSeller.Email).Return(&testSeller, nil).Times(2)

	t.Run("Testing for Bad/Wrong Request", func(t *testing.T) {
		mockDB.EXPECT().FindIndividualSellerShop(testShopID).Return(&testSeller, nil)
		mockDB.EXPECT().FindSellerIndividualProduct(sellerID).Return(&indProduct, nil)
		mockDB.EXPECT().FindCartProductSeller(sellerID, Id).Return(nil,
			fmt.Errorf("an error occurred"))
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/seller/remaining/product/count",
			strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Error", "not", "find", "cart")
	})

	t.Run("Testing for Successful Request", func(t *testing.T) {
		mockDB.EXPECT().FindIndividualSellerShop(testShopID).Return(&testSeller, nil)
		mockDB.EXPECT().FindSellerIndividualProduct(sellerID).Return(&indProduct, nil)
		mockDB.EXPECT().FindCartProductSeller(sellerID, Id).Return(&testCartProductTwo, nil)
		//mockDB.EXPECT().UpdateSellerProductQuantity(sellerID, productID, remainingQuantity).Return(&indProduct, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/seller/remaining/product/count",
			strings.NewReader(string(bodyJSON)))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "Seller", "Remaining", "Product")
	})

}
