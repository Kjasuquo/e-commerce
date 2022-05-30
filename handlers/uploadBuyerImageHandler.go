package handlers

import (
	"github.com/decadevs/shoparena/models"
	"github.com/decadevs/shoparena/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// UploadImageHandler uploads a buyer's profile picture
func (h *Handler) UploadBuyerImageHandler(c *gin.Context) {
	if userI, exists := c.Get("user"); exists {
		if user, ok := userI.(*models.Buyer); ok {

			const maxSize = int64(2048000) // allow only 2MB of file size
			r := c.Request
			err := r.ParseMultipartForm(maxSize)
			if err != nil {
				log.Printf("parse image error: %v\n", err)
				c.JSON(http.StatusBadRequest, []string{"image too large"})
				return
			}
			file, fileHeader, err := r.FormFile("profile_picture")
			if err != nil {
				log.Println("error getting profile picture", err)
				c.JSON(http.StatusBadRequest, []string{"image not supplied"})
				return
			}
			defer file.Close()
			fileExtension, ok := services.CheckSupportedFile(strings.ToLower(fileHeader.Filename))
			log.Println(filepath.Ext(strings.ToLower(fileHeader.Filename)))
			if ok {
				log.Println(fileExtension)
				c.JSON(http.StatusBadRequest, []string{fileExtension})
				return
			}
			session, tempFileName, err := services.PreAWS(fileExtension, "profile_picture")
			if err != nil {
				log.Printf("could not upload file: %v\n", err)
			}
			url, err := h.DB.UploadFileToS3(session, file, tempFileName, fileHeader.Size)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, []string{"an error occurred while uploading the image"})
				return
			}
			buyerID := user.ID
			user.Image = url
			err = h.DB.UpdateBuyerImageURL(user.Username, user.Image, buyerID)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, []string{"an error occurred while uploading the image"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"imageurl": user.Image,
			})
			return

		}
	}
	c.JSON(http.StatusUnauthorized, []string{"unable to retrieve authenticated user"})
}
