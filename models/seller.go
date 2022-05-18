package models

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	User
	Product []Product `json:"product" gorm:"oneToMany"`
	Orders  []Order   `json:"orders" gorm:"oneToMany"`
	Rating  int       `json:"rating"`
}
