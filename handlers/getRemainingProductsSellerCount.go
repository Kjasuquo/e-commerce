package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) GetRemainingProductsCountSellerCount(c *gin.Context) {
	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	sellerID := seller.ID

	//get that the particular seller shop
	seller, err = h.DB.FindIndividualSellerShop(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist ; Seller with this ID not found in product table",
			"error":   err.Error(),
		})
		return
	}

	//get seller slice of products
	sellerProducts := seller.Product

	//number of products in slice of products
	productLength := len(sellerProducts)

	var (
		countAllProduct       int
		soldProductCount      int
		totalQuantity         uint
		sellerProductID       uint
		remainingProductCount int
	)

	//total count based on the category
	for i := 0; i < productLength; i++ {
		countAllProduct++
		totalQuantity += seller.Product[i].Quantity
	}

	indProduct, err := h.DB.FindSellerIndividualProduct(sellerID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"Message": "Error Exist ; Product not found in product table",
			"error":   err.Error(),
		})
	}

	productID := indProduct.ID

	cartProduct, err := h.DB.FindCartProductSeller(sellerID, productID)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist, could not find cart",
		})
		return
	}

	for i := 0; i < productLength; i++ {
		if cartProduct.OrderStatus {
			if cartProduct.TotalQuantity > 0 {
				soldProductCount++
				remainingProductCount += int(totalQuantity - cartProduct.TotalQuantity)
			}
		}
	}

	//get the remaining product count
	finalRmProductCount := uint(remainingProductCount)

	for i := 0; i < productLength; i++ {
		sellerProductID += seller.Product[i].Quantity
	}

	previousQuantity := indProduct.Quantity
	newQuantityLeft := previousQuantity - uint(soldProductCount)

	c.IndentedJSON(http.StatusOK, gin.H{
		"Message":         "Seller Remaining Product",
		"Total_Product":   totalQuantity,
		"Total_Sold":      soldProductCount,
		"Total_Remaining": finalRmProductCount,
		"new_quantity":    newQuantityLeft,
	})
}
