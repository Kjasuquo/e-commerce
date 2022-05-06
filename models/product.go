package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ShopName        string `json:"shop_name"`
	ProductName     string `json:"product_name"`
	ProductPrice    uint   `json:"product_price"`
	ProductCategory string `json:"product_category"`
	ProductImage    string `json:"product_image"`
	ProductDetails  string `json:"product_details"`
	Rating          uint   `json:"rating"`
	Quantity        uint   `json:"quantity"`
}
