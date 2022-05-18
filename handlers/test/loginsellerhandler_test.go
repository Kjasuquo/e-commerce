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
)

func TestSellerLogin(t *testing.T) {

	//import new controller
	ctxl := gomock.NewController(t)
	defer ctxl.Finish()

	mockdb := mock_database.NewMockDB(ctxl)

	mockhandler := &handlers.Handler{
		DB: mockdb,
	}

	rout, _ := router.SetupRouter(mockhandler)
	t.Run("testing bad request", func(t *testing.T) {
		usermk := &struct {
			Email string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email: "",
			Password: "12345566666",
		}
		bytes, _ := json.Marshal(usermk)
		ht := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/loginseller", strings.NewReader(string(bytes)))

		rout.ServeHTTP(ht, req)
		assert.Equal(t, http.StatusBadRequest, ht.Code)
		assert.Contains(t, ht.Body.String(), "bad request")
	})
	t.Run("find seller by email", func(t *testing.T) {
		hash, _ := handlers.HashPassword("12345566666")
		seller := &models.Seller{
			User: models.User{
				Email:     "mike123",
				PasswordHash: hash,
			},
		}
		usermk := &struct {
			Email string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}{
			Email: "mike123",
			Password: "12345566666",
		}
		mockdb.EXPECT().FindSellerByEmail(seller.Email).Return(seller, nil)
		bytes, _ := json.Marshal(usermk)
		ht := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/loginseller", strings.NewReader(string(bytes)))

		rout.ServeHTTP(ht, req)
		assert.Equal(t, http.StatusOK, ht.Code)
		fmt.Println(ht.Body.String())
		assert.Contains(t, ht.Body.String(), usermk.Email)
	})

}
