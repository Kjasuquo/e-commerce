package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	mock_database "github.com/decadevs/shoparena/database/mocks"
	"github.com/decadevs/shoparena/handlers"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/router"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetProductById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mock_database.NewMockDB(ctrl)

	h := &handlers.Handler{DB: mockDB}

	route, _ := router.SetupRouter(h)

	testGorm := gorm.Model{
		ID:        1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	category := models.Category{
		testGorm,
		"Handgloves",
	}
	images := models.Image{
		testGorm, 5, "gre.com",
	}
	sliceImages := []models.Image{images}

	//an instance of product
	product := models.Product{
		Model:       testGorm,
		SellerId:    2,
		CategoryId:  4,
		Category:    category,
		Title:       "RubberGloves",
		Description: "Blue, latex",
		Price:       600,
		Images:      sliceImages,
		Rating:      6,
		Quantity:    1000,

	}
	bodyJSON, err := json.Marshal(product)
	if err != nil {
		t.Fail()
	}
	t.Run("Testing For Error", func(t *testing.T) {
		id := uint(1)
		mockDB.EXPECT().GetProductByID(id).Return(nil, errors.New("Error Exist"))
		rw := httptest.NewRecorder()
		idVal := strconv.Itoa(int(testGorm.ID))
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/"+idVal, strings.NewReader(string(bodyJSON)))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "Error Exist in ")
	})
	t.Run("getting product by ID successful", func(t *testing.T) {
		id := 1
		mockDB.EXPECT().GetProductByID(uint(id)).Return(&product, nil)
		rw := httptest.NewRecorder()
		idVal := strconv.Itoa(int(testGorm.ID))
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/product/"+idVal, strings.NewReader(string(bodyJSON)))
		route.ServeHTTP(rw, req)
		fmt.Println(rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
		assert.Contains(t, rw.Body.String(), "product gotten by ID")
	})
}
