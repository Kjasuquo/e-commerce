package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/mail"
)

type User struct {
	gorm.Model
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email" gorm:"unique"`
	Username        string `json:"username" gorm:"unique"`
	Password        string `json:"password,omitempty" gorm:"-"`
	ConfirmPassword string `json:"confirm_password,omitempty" gorm:"-"`
	PasswordHash    string `json:"-"`
	Address         string `json:"address"`
	PhoneNumber     string `json:"phone_number"`
	Image           string `json:"image"`
	IsActive        bool   `json:"status"`
	Token           string `json:"token"`
}

type UpdateUser struct {
	FirstName   string `json:"first_name" binding:"required" form:"first_name"`
	LastName    string `json:"last_name" binding:"required" form:"last_name"`
	PhoneNumber string `json:"phone_number" binding:"required" form:"phone_number"`
	Email       string `json:"email" binding:"required,email" form:"email"`
	Address     string `json:"address"  form:"address"`
	Rating      uint   `json:"rating"`
}

type UpdateRating struct {
	Rating                  uint `json:"rating"`
	TotalRatings            uint `json:"total_ratings"`
	NumberOfRatingsReceived uint `json:"number_of_ratings_received"`
}

func (user *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return nil
}

func (user *User) ValidMailAddress() bool {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return false
	}
	return true
}
