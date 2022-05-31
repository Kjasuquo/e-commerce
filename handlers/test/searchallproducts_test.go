package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSearchAllProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)
	mockMail := mock_database.NewMockMailer(ctrl)

	h := &handlers.Handler{DB: mockDB, Mail: mockMail}

	route, _ := router.SetupRouter(h)

	category1 := models.Category{
		Model: gorm.Model{ID: 1},
		Name:  "shirt",
	}

	category2 := models.Category{
		Model: gorm.Model{ID: 2},
		Name:  "short",
	}

	category3 := models.Category{
		Model: gorm.Model{ID: 3},
		Name:  "shoe",
	}

	product := []models.Product{
		{
			SellerId:    1,
			CategoryId:  category2.ID,
			Title:       "combat shorts",
			Description: "nice shorts",
			Price:       5000,
			Rating:      10,
			Quantity:    8,
		},
		{
			SellerId:    2,
			CategoryId:  category2.ID,
			Title:       "beach shorts",
			Description: "for beach outing",
			Price:       4000,
			Rating:      9,
			Quantity:    15,
		},
		{
			SellerId:    3,
			CategoryId:  category1.ID,
			Title:       "vintage shirts",
			Description: "cool for casual wears",
			Price:       4500,
			Rating:      10,
			Quantity:    7,
		},

		{
			SellerId:    1,
			CategoryId:  category1.ID,
			Title:       "combat shirts",
			Description: "ready for war",
			Price:       7600,
			Rating:      8,
			Quantity:    9,
		},
		{
			SellerId:    2,
			CategoryId:  category1.ID,
			Title:       "beach shirts",
			Description: "for beach outing",
			Price:       3000,
			Rating:      9,
			Quantity:    15,
		},
		{
			SellerId:    3,
			CategoryId:  category2.ID,
			Title:       "bum shorts",
			Description: "sexy bum shorts",
			Price:       2500,
			Rating:      9,
			Quantity:    7,
		},
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		t.Fail()
	}

	productCat1 := []models.Product{

		{
			SellerId:    3,
			CategoryId:  category1.ID,
			Title:       "vintage shirts",
			Description: "cool for casual wears",
			Price:       4500,
			Rating:      10,
			Quantity:    7,
		},

		{
			SellerId:    1,
			CategoryId:  category1.ID,
			Title:       "combat shirts",
			Description: "ready for war",
			Price:       7600,
			Rating:      8,
			Quantity:    9,
		},
		{
			SellerId:    2,
			CategoryId:  category1.ID,
			Title:       "beach shirts",
			Description: "for beach outing",
			Price:       3000,
			Rating:      9,
			Quantity:    15,
		},
	}

	productcat1JSON, err := json.Marshal(productCat1)
	if err != nil {
		t.Fail()
	}

	t.Run("Testing for empty queries", func(t *testing.T) {
		mockDB.EXPECT().SearchProduct("", "", "", "").Return(product, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/searchproducts", strings.NewReader(string(productJSON)))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Code)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(productJSON))

	})

	t.Run("Testing for error", func(t *testing.T) {
		mockDB.EXPECT().SearchProduct("", "", category3.Name, "").Return(nil, err)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/searchproducts?category=shoe", strings.NewReader(string(productJSON)))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Code)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "no such product")

	})

	t.Run("Testing for queries", func(t *testing.T) {
		mockDB.EXPECT().SearchProduct("3000", "8000", "shirt", "shirts").Return(productCat1, nil)
		rw := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/searchproducts?category=shirt&lower-price=3000&upper-price=8000&name=shirts", strings.NewReader(string(productcat1JSON)))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Code)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(productcat1JSON))

	})
}
