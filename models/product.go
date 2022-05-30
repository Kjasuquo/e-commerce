package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	SellerId uint `json:"seller_id"`

	CategoryId  uint `gorm:"foreignKey:categories(id)" json:"category_id"`
	Category    Category
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       uint    `json:"price"`
	Images      []Image `json:"images" gorm:"oneToMany"`
	Rating      uint    `json:"rating"`
	Quantity    uint    `json:"quantity"`
}
