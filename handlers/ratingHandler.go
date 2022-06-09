package handlers

import (
	"github.com/decadevs/shoparena/models"
	"github.com/gin-gonic/gin"
	"log"
)

func (h *Handler) SellerRating(c *gin.Context) {
	var rating models.RatingRequest
	userI, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{"message": "internal server get error"})
	}
	_, ok := userI.(*models.Buyer)
	if !ok {
		c.JSON(500, gin.H{"message": "invalid assert"})
	}
	err := c.BindJSON(&rating)
	if err != nil {
		c.JSON(404, gin.H{"message": "please fill al fields"})
		c.Abort()
		return
	}
	seller, err := h.DB.FindSellerById(rating.Id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Something Went Wrong!"})
		c.Abort()
		return
	}
	seller.TotalRatings += rating.Rating
	seller.NumberOfRatingsReceived += 1
	log.Println("test1")

	seller.Rating = seller.TotalRatings / seller.NumberOfRatingsReceived
	var updateDetail = models.UpdateRating{
		Rating:                  seller.Rating,
		TotalRatings:            seller.TotalRatings,
		NumberOfRatingsReceived: seller.NumberOfRatingsReceived,
	}
	err = h.DB.UpdateSellerRating(seller.ID, &updateDetail)
	if err != nil {
		c.JSON(500, gin.H{"message": "Something went wrong, please try again"})
		c.Abort()
		return
	}
	log.Println(seller.Rating)
	c.JSON(200, gin.H{"message": "thank you for your feedback"})
}
func (h *Handler) ProductRating(c *gin.Context) {
	var ratingRequest models.RatingRequest
	userI, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{"message": "internal server get  error"})
	}
	_, ok := userI.(*models.Buyer)
	if !ok {
		c.JSON(500, gin.H{"message": "invalid assert"})
	}
	err := c.BindJSON(&ratingRequest)
	if err != nil {
		c.JSON(404, gin.H{"message": "please fill al fields"})
		c.Abort()
		return
	}
	product, err := h.DB.FindProductById(ratingRequest.Id)
	if err != nil {
		c.JSON(500, gin.H{"message": "Something Went Wrong!"})
		c.Abort()
		return
	}
	product.TotalRatings += ratingRequest.Rating
	product.NumberOfRatingsReceived += 1
	log.Println("test1")

	product.Rating = product.TotalRatings / product.NumberOfRatingsReceived
	var updateDetail = models.UpdateRating{
		Rating:                  product.Rating,
		TotalRatings:            product.TotalRatings,
		NumberOfRatingsReceived: product.NumberOfRatingsReceived,
	}
	err = h.DB.UpdateProductRating(product.ID, &updateDetail)
	if err != nil {
		c.JSON(500, gin.H{"message": "Something went wrong, please try again"})
		c.Abort()
		return
	}
	log.Println(product.Rating)
	c.JSON(200, gin.H{"message": "thank you for your feedback"})
}
