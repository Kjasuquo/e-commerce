package models

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	User
	Product                 []Product `json:"product" gorm:"oneToMany"`
	Orders                  []Order   `json:"orders" gorm:"oneToMany"`
	Rating                  uint      `json:"rating"`
	TotalRatings            uint      `json:"total_ratings"`
	NumberOfRatingsReceived uint      `json:"number_of_ratings_received"`
}
