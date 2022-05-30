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
	"log"
	"net/http"
	"os"

	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateProduct(t *testing.T) {
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

		product := models.Product{

			SellerId:    1,
			CategoryId:  2,
			Title:       "plumbing",
			Description: "Building",
			Price:       6700,
			Rating:      5,
			Quantity:    3,
		}

		bodyJSON, err := json.Marshal(product)
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		mockDB.EXPECT().TokenInBlacklist(gomock.Any()).Return(false).Times(2)
		mockDB.EXPECT().FindSellerByEmail(seller.Email).Return(seller, nil).Times(2)

		t.Run("testing for created update", func(t *testing.T) {
			mockDB.EXPECT().CreateProduct(product).Return(nil).AnyTimes()
			rw := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/api/v1/createproduct/1",
				strings.NewReader(string(bodyJSON)))
			if err != nil {
				fmt.Printf("error occured %v \n", err)
				return
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *acc))
			route.ServeHTTP(rw, req)
			assert.Equal(t, http.StatusInternalServerError, rw.Code)
			assert.Contains(t, rw.Body.String(), "error adding product")
			t.Log(err)
			route.ServeHTTP(rw, req)
			assert.NotEqual(t, http.StatusCreated, rw.Code)
			assert.Contains(t, rw.Body.String(), "product successfully added")
		})
	}
}
