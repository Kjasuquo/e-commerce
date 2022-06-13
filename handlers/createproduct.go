package handlers

import (
	"fmt"
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/server/response"
	"github.com/decadevs/shoparena/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func (h *Handler) CreateProducts(c *gin.Context) {
	//get the user id from a logged-in user
	user, exists := c.Get("user")
	if !exists {
		log.Printf("can't get user from context\n")
		response.JSON(c, "", http.StatusUnauthorized, nil, []string{"you are not logged in"})
		return
	}
	seller := user.(*models.Seller)
	//userId := userI.(*models.User).ID

	form, err := c.MultipartForm()

	if err != nil {
		log.Printf("error parsing multipart form: %v", err)
		response.JSON(c, "", http.StatusInternalServerError, nil, []string{"internal server error"})
		return
	}

	formImages := form.File["images"]
	images := []models.Image{}
	log.Println(formImages)
	log.Println(images)

	// upload the images to aws.
	for _, f := range formImages {
		file, err := f.Open()
		if err != nil {

		}
		fileExtension, ok := services.CheckSupportedFile(strings.ToLower(f.Filename))
		log.Printf(filepath.Ext(strings.ToLower(f.Filename)))
		fmt.Println(fileExtension)
		if ok {
			log.Println(fileExtension)
			response.JSON(c, "", http.StatusBadRequest, nil, []string{fileExtension + " image file type is not supported"})
			return
		}

		session, tempFileName, err := services.PreAWS(fileExtension, "product")
		if err != nil {
			log.Println("could not upload file", err)
		}

		url, err := h.DB.UploadFileToS3(session, file, tempFileName, f.Size)
		if err != nil {
			log.Println(err)
			response.JSON(c, "", http.StatusInternalServerError, nil, []string{"an error occurred while uploading the image"})
			return
		}

		log.Printf("filename: %v", f.Filename)

		img := models.Image{
			Url: url,
		}
		images = append(images, img)
	}

	price, err := strconv.Atoi(c.PostForm("price"))
	if err != nil {
		log.Println(err)
		response.JSON(c, "", http.StatusBadRequest, "in price", []string{err.Error()})
		return
	}

	// err := strconv.Atoi(c.PostForm("no_of_ratings"))
	rating, err := strconv.Atoi(c.PostForm("rating"))
	if err != nil {
		response.JSON(c, "", http.StatusBadRequest, "in rating", []string{err.Error()})
		return
	}

	quantity, err := strconv.Atoi(c.PostForm("quantity"))
	if err != nil {
		log.Println(err)

		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return

		response.JSON(c, "", http.StatusBadRequest, "in quantity", []string{err.Error()})

	}
	CategoryID, err := strconv.Atoi(c.PostForm("category_id"))
	if err != nil {
		log.Println(err)

		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
    
		response.JSON(c, "", http.StatusBadRequest, "in category_id", []string{err.Error()})

	}

	products := models.Product{
		Category: models.Category{
			Name: c.PostForm("name"),
			Model: gorm.Model{
				ID: uint(CategoryID),
			},
		},
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Price:       uint(price),
		SellerId:    seller.ID,
		Images:      images,
		Rating:      uint(rating),
		Quantity:    uint(quantity),
	}
	log.Println(products, CategoryID)
	err = h.DB.CreateProduct(products)
	if err != nil {
		response.JSON(c, "", http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

	log.Println(products.Images)
	c.JSON(http.StatusOK, products)

}
