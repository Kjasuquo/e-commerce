package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

//GET TOTAL PRODUCT COUNT FOR A SELLER
func (h *Handler) GetTotalProductCountForSeller(c *gin.Context) {

	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	sellerID := strconv.Itoa(int(seller.ID))

	//find seller with the retrieved ID and return the seller and its preloaded product
	Seller, err := h.DB.FindIndividualSellerShop(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist ; Seller with this ID not found in product table",
			"error":   err.Error(),
		})
		return
	}

	//find seller product , retrived by ID and return the seller product quantity
	product, err := h.DB.FindSellerProduct(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist count not find seller product",
			"error":   err.Error(),
		})
		return
	}

	//get seller slice of products
	sellerProducts := Seller.Product

	//number of products in slice of products
	productLength := len(sellerProducts)

	//count all products in slice == length pf slice
	var countAllProduct int

	//sum of the quantity of each product
	var total_quantity uint

	//if seller do not have product return, i.e. do run the below code
	if productLength == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"Message": "Seller do not have product to sell",
			"Value":   countAllProduct,
		})
		return
	}
	//total count based on the category
	for i := 0; i < productLength; i++ {
		countAllProduct++
	}

	//total count based on the individual product
	for i := 0; i < productLength; i++ {
		total_quantity += product[i].Quantity
	}

	//individual quantity of products can be gotten by
	quantity := product[1].Quantity

	//total count based on an individual product
	c.IndentedJSON(http.StatusOK, gin.H{
		"Message":                "Seller has " + strconv.Itoa(countAllProduct) + " different categories of product to sell",
		"Seller_ID":              sellerID,
		"Product_Count":          countAllProduct,
		"Product_Quantity":       total_quantity,
		"First_Product_Quantity": quantity,
	})

}

func (h *Handler) SellerIndividualProduct(c *gin.Context) {
	//authorize and authenticate seller
	seller, err := h.GetUserFromContext(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, []string{"internal server error"})
		return
	}
	sellerID := strconv.Itoa(int(seller.ID))

	product, err := h.DB.FindSellerProduct(sellerID)
	if err != nil {
		log.Println("Error finding information in database:", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "Error Exist Count not find individual seller product",
			"error":   err.Error(),
		})
		return
	}

	if len(product) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Message": "No Product Found In Shop",
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Message": "Seller Has Product In Shop",
		"Product": product,
	})

}
