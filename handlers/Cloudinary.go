package handlers

//
//import (
//	"github.com/cloudinary/cloudinary-go/v2"
//	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
//	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//	"strings"
//)
//
//func UploadFile(c *gin.Context) {
//	cld, _ := cloudinary.NewFromURL("cloudinary://865474734328532:lSqgNEqmjihsvOiJFFUnQHoxLI8@kjasuquo")
//	// Get the preferred name of the file if its not supplied
//	log.Println("HEY 1")
//	fileName := c.PostForm("image")
//	log.Println("HEY 2")
//	// Add tags
//	fileTags := c.PostForm("tags")
//	log.Println("HEY 3")
//	file, _, err := c.Request.FormFile("file")
//	log.Println("HEY 4")
//	if err != nil {
//		log.Println("HEY 5 ", err)
//		c.JSON(http.StatusBadRequest, gin.H{
//			"error":   err,
//			"message": "Failed to upload",
//		})
//	}
//
//	result, err := cld.Upload.Upload(c, file, uploader.UploadParams{
//		PublicID: fileName,
//		// Split the tags by comma
//		Tags: strings.Split(",", fileTags),
//	})
//
//	if err != nil {
//		c.String(http.StatusConflict, "Upload to cloudinary failed")
//	}
//
//	c.JSON(http.StatusCreated, gin.H{
//		"message": "Successfully uploaded the file",
//		//"secureURL": result.SecureURL,
//		"publicURL": result,
//	})
//}
//
////func getUploadedFiles(c *gin.Context) {
////	// Create our instance
////	cld, _ := cloudinary.NewFromURL("your_cloud_url")
////	fileName := c.Param("accessId")
////
////	// Access the filename using a desired file access id.
////	result, err := cld.Admin.Asset(c, admin.AssetParams{
////		AccessID: fileName,
////	})
////	if err != nil {
////		c.String(http.StatusNotFound, "We were unable to find the file requested")
////	}
////
////	c.JSON(http.StatusAccepted, gin.H{
////		"message":    "Successfully found the image",
////		"secureURL":  result.SecureURL,
////		"publicURL":  result.URL,
////		"created_at": result.CreatedAt.String(),
////	})
////}
//func GetUploadedFilesForMultipleImages(c *gin.Context) {
//	cld, _ := cloudinary.NewFromURL("your_cloud_url")
//	var urls []string
//
//	searchQ := search.Query{
//		Expression: "resource_type:image AND uploaded_at>1d AND bytes<1m",
//		SortBy: []search.SortByField{
//			{
//				"created_at": search.Descending,
//			},
//		},
//		MaxResults: 10,
//	}
//
//	results, err := cld.Admin.Search(c, searchQ)
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{
//			"error":   err,
//			"message": "Failed to query your files",
//		})
//	}
//
//	for _, asset := range results.Assets {
//		urls = append(urls, asset.SecureURL)
//	}
//
//	c.JSON(http.StatusAccepted, gin.H{
//		"success": true,
//		"data":    urls,
//	})
//}
//
//func UpdateFile(c *gin.Context) {
//	cld, _ := cloudinary.NewFromURL("your_url")
//	fileId := c.Param("publicId")
//	newFileName := c.PostForm("fileName")
//
//	// Access the filename using a desired filename
//	result, err := cld.Upload.Rename(c, uploader.RenameParams{
//		FromPublicID: fileId,
//		ToPublicID:   newFileName,
//	})
//	if err != nil {
//		c.String(http.StatusNotFound, "We were unable to find the file requested")
//	}
//
//	c.JSON(http.StatusAccepted, gin.H{
//		"message":    "Successfully found the image",
//		"secureURL":  result.SecureURL,
//		"publicURL":  result.URL,
//		"created_at": result.CreatedAt.String(),
//	})
//}
//
//func DeleteFile(c *gin.Context) {
//	cld, _ := cloudinary.NewFromURL("cloudinary://865474734328532:lSqgNEqmjihsvOiJFFUnQHoxLI8@kjasuquo")
//	fileId := c.Param("assetId")
//	result, err := cld.Upload.Destroy(c, uploader.DestroyParams{
//		PublicID: fileId,
//	})
//
//	if err != nil {
//		c.String(http.StatusBadRequest, "File could not be deleted")
//	}
//
//	c.JSON(http.StatusContinue, result.Result)
//}
