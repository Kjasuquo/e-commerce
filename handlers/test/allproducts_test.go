package test

import (
	"encoding/json"
	"fmt"
	mockdatabase "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mockdatabase.NewMockDB(ctrl)
	h := &handlers.Handler{
		DB: mockDB,
	}
	route, _ := router.SetupRouter(h)

	category := models.Category{
		Model: gorm.Model{ID: 1},
		Name:  "Shirt",
	}

	product := []models.Product{
		{
			Model:       gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
			SellerId:    2,
			CategoryId:  category.ID,
			Title:       "household",
			Description: "home items",
			Price:       78,
			Rating:      8,
			Quantity:    4903,
		},
		{
			Model:       gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
			SellerId:    2,
			CategoryId:  3,
			Title:       "household",
			Description: "home items",
			Price:       78,
			Rating:      8,
			Quantity:    4903,
		},
	}

	bodyJSON, err := json.Marshal(product)
	if err != nil {
		t.Fail()
	}

	t.Run("Testing for all product", func(t *testing.T) {
		mockDB.EXPECT().GetAllProducts().Return(product)
		rw := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet,
			"/api/v1/products",
			strings.NewReader(string(bodyJSON)))
		if err != nil {
			fmt.Printf("errrr here %v \n", err)
			return
		}
		route.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), string(bodyJSON))
	})

}
