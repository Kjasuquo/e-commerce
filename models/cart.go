package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	BuyerID uint          `json:"buyers_id"`
	Product []CartProduct `json:"product"`
}

type CartProduct struct {
	gorm.Model
	CartID        uint `json:"cart_id"`
	ProductID     uint `json:"product_id"`
	TotalPrice    uint `json:"total_price"`
	TotalQuantity uint `json:"total_quantity"`
	OrderStatus   bool `json:"order_status"`
	SellerId      uint `json:"seller_id"`
	BuyerId       uint `json:"buyer_id"`
}

type ProductDetails struct {
	gorm.Model
	Name          string
	Price         uint
	Quantity      uint
	Images        []Image
	CartProductID uint
}
