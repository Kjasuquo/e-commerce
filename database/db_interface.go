package database

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/decadevs/shoparena/models"
	"github.com/joho/godotenv"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// DB provides access to the different db

type DB interface {
	CreateBuyer(user *models.Buyer) (*models.Buyer, error)
	CreateSeller(user *models.Seller) (*models.Seller, error)
	FindAllSellersExcept(except string) ([]models.Seller, error)
	FindBuyerByEmail(email string) (*models.Buyer, error)
	FindBuyerByPhone(phone string) (*models.Buyer, error)
	FindBuyerByUsername(username string) (*models.Buyer, error)
	FindSellerByEmail(email string) (*models.Seller, error)
	FindSellerByPhone(phone string) (*models.Seller, error)
	UpdateUserImageURL(username, url string) error
	FindSellerByUsername(username string) (*models.Seller, error)
	SearchProduct(lowerPrice, upperPrice, category, name string) ([]models.Product, error)
	TokenInBlacklist(token *string) bool
	UpdateBuyerProfile(id uint, update *models.UpdateUser) error
	UpdateSellerProfile(id uint, update *models.UpdateUser) error
	UploadFileToS3(h *session.Session, file multipart.File, fileName string, size int64) (string, error)
	BuyerUpdatePassword(password, newPassword string) (*models.Buyer, error)
	SellerUpdatePassword(password, newPassword string) (*models.Seller, error)
	BuyerResetPassword(email, newPassword string) (*models.Buyer, error)
	CreateBuyerCart(cart *models.Cart) (*models.Cart, error)
	FindIndividualSellerShop(sellerID string) (*models.Seller, error)
	GetAllProducts() []models.Product
	UpdateProductByID(Id uint, prod models.Product) error
	GetAllSellers() ([]models.Seller, error)
	GetProductByID(id uint) (*models.Product, error)
	FindSellerProduct(sellerID string) ([]models.Product, error)
	GetAllBuyerOrder(buyerId uint) ([]models.Order, error)
	GetAllSellerOrder(sellerId uint) ([]models.Order, error)
	GetAllSellerOrderCount(sellerId uint) (int, error)
	FindPaidProduct(sellerID string) ([]models.CartProduct, error)
	AddToCart(product models.Product, buyer *models.Buyer) error
	GetCartProducts(buyer *models.Buyer) ([]models.CartProduct, error)
	ViewCartProducts(addedProducts []models.CartProduct) ([]models.ProductDetails, error)
	DeletePaidFromCart(cartID uint) error
	GetSellersProducts(sellerID uint) ([]models.Product, error)
}

// Mailer interface to implement mailing service
type Mailer interface {
	SendMail(subject, body, to, Private, Domain string) error
	GenerateNonAuthToken(UserEmail string, secret string) (*string, error)
	DecodeToken(token, secret string) (string, error)
}

//Paystack interface
type Paystack interface {
	InitializePayment(info []byte) (string, error)
	Callback(reference string) (*http.Response, error)
}

// ValidationError defines error that occur due to validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}

type DBParams struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
}

func InitDBParams() DBParams {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("PDB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PDB_PORT")

	return DBParams{
		Host:     host,
		User:     user,
		Password: password,
		DbName:   dbName,
		Port:     port,
	}
}
