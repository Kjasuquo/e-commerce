package database

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/decadevs/shoparena/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

//PostgresDb implements the DB interface
type PostgresDb struct {
	DB *gorm.DB
}

// Init sets up the mongodb instance
func (pdb *PostgresDb) Init(host, user, password, dbName, port string) error {
	fmt.Println("connecting to Database.....")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos",
		host, user, password, dbName, port)
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if db == nil {
		return fmt.Errorf("database was not initialized")
	} else {
		fmt.Println("Connected to Database")
	}

	pdb.DB = db
	err = pdb.PrePopulateTables()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}

func (pdb *PostgresDb) PrePopulateTables() error {
	err := pdb.DB.AutoMigrate(&models.Category{}, &models.Seller{}, &models.Product{}, &models.Image{},
		&models.Buyer{}, &models.Cart{}, &models.CartProduct{}, &models.Order{}, &models.Blacklist{})
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}
	categories := []models.Category{{Name: "fashion"}, {Name: "electronics"}, {Name: "health & beauty"}, {Name: "baby products"}, {Name: "phones & tablets"}, {Name: "food drinks"}, {Name: "computing"}, {Name: "sporting goods"}, {Name: "others"}}
	result := pdb.DB.Find(&models.Category{})
	if result.RowsAffected < 1 {
		pdb.DB.Create(&categories)
	}

	user := models.User{
		Model:           gorm.Model{},
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "jdoe@gmail.com",
		Username:        "JD Baba",
		Password:        "12345678",
		ConfirmPassword: "12345678",
		PasswordHash:    "$2a$12$T2wSf1qgpTyhLOons3u4JOCqCwKDDL4J3UhGdOTEBL/CmAS/RNCPm",
		Address:         "aso rock",
		PhoneNumber:     "09091919292",
		Image:           "https://i.ibb.co/5jwDfyF/Photo-on-24-11-2021-at-20-45.jpg",
		IsActive:        true,
		Token:           "",
	}
	buyer := models.Buyer{
		Model:  gorm.Model{},
		User:   user,
		Orders: nil,
	}
	result = pdb.DB.Where("buyer = ?", "John").Find(&buyer)

	if result.RowsAffected < 1 {
		pdb.DB.Create(&buyer)
	}

	seller := models.Seller{
		Model:   gorm.Model{},
		User:    user,
		Product: nil,
		Orders:  nil,
		Rating:  5,
	}
	result = pdb.DB.Where("seller = ?", "John").Find(&seller)

	if result.RowsAffected < 1 {
		pdb.DB.Create(&seller)
	}
	return nil
}

//GET ALL PRODUCTS FROM DB
func (pdb *PostgresDb) GetAllProducts() []models.Product {
	var products []models.Product
	if err := pdb.DB.Find(&products).Error; err != nil {
		log.Println("Could not find product", err)
	}
	return products
}

//UPDATE PRODUCT BY ID
func (pdb *PostgresDb) UpdateProductByID(Id uint, prod models.Product) error {
	products := models.Product{}

	err := pdb.DB.Model(&products).Where("id = ?", Id).Update("title", prod.Title).
		Update("description", prod.Description).Update("price", prod.Price).
		Update("rating", prod.Rating).Update("quantity", prod.Quantity).Error
	if err != nil {
		fmt.Println("error in updating in postgres db")
		return err
	}
	return nil
}

// SearchProduct Searches all products from DB
func (pdb *PostgresDb) SearchProduct(lowerPrice, upperPrice, categoryName, name string) ([]models.Product, error) {
	categories := models.Category{}
	var products []models.Product

	LPInt, _ := strconv.Atoi(lowerPrice)
	UPInt, _ := strconv.Atoi(upperPrice)

	if categoryName == "" {
		err := pdb.DB.Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return products, nil
	} else {
		err := pdb.DB.Where("name = ?", categoryName).First(&categories).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	category := categories.ID

	if LPInt == 0 && UPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price <= ?", uint(UPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if UPInt == 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price >= ?", uint(LPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt != 0 && UPInt != 0 && name == "" {
		err := pdb.DB.Where("category_id = ?", category).Where("price >= ?", uint(LPInt)).
			Where("price <= ?", uint(UPInt)).Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && UPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if LPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price <= ?", uint(UPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else if UPInt == 0 && name != "" {
		err := pdb.DB.Where("category_id = ?", category).
			Where("price >= ?", uint(LPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		err := pdb.DB.Where("category_id = ?", category).Where("price >= ?", uint(LPInt)).
			Where("price <= ?", uint(UPInt)).
			Where("title LIKE ?", "%"+name+"%").Find(&products).Error
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return products, nil
}

// CreateSeller creates a new Seller in the DB
func (pdb *PostgresDb) CreateSeller(user *models.Seller) (*models.Seller, error) {
	var err error
	user.CreatedAt = time.Now()
	user.IsActive = true
	err = pdb.DB.Create(user).Error
	return user, err
}

// CreateBuyer creates a new Buyer in the DB
func (pdb *PostgresDb) CreateBuyer(user *models.Buyer) (*models.Buyer, error) {
	var err error
	user.CreatedAt = time.Now()
	user.IsActive = true
	err = pdb.DB.Create(user).Error
	return user, err
}

//CreateBuyerCart creates a new cart for the buyer
func (pdb *PostgresDb) CreateBuyerCart(cart *models.Cart) (*models.Cart, error) {
	var err error
	cart.CreatedAt = time.Now()
	err = pdb.DB.Create(cart).Error
	return cart, err
}

// FindSellerByUsername finds a user by the username
func (pdb *PostgresDb) FindSellerByUsername(username string) (*models.Seller, error) {
	user := &models.Seller{}

	if err := pdb.DB.Where("username = ?", username).First(user).Error; err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("user inactive")
	}
	return user, nil
}

// FindBuyerByUsername finds a user by the username
func (pdb *PostgresDb) FindBuyerByUsername(username string) (*models.Buyer, error) {
	buyer := &models.Buyer{}

	if err := pdb.DB.Where("username = ?", username).First(buyer).Error; err != nil {
		return nil, err
	}
	if !buyer.IsActive {
		return nil, errors.New("user inactive")
	}
	return buyer, nil
}

// FindSellerByEmail finds a user by email
func (pdb *PostgresDb) FindSellerByEmail(email string) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := pdb.DB.Where("email = ?", email).First(seller).Error; err != nil {
		return nil, errors.New(email + " does not exist" + " seller not found")
	}

	return seller, nil
}

// FindBuyerByEmail finds a user by email
func (pdb *PostgresDb) FindBuyerByEmail(email string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Where("email = ?", email).First(buyer).Error; err != nil {
		return nil, errors.New(email + " does not exist" + " buyer not found")
	}

	return buyer, nil
}

// FindSellerByPhone finds a user by the phone
func (pdb PostgresDb) FindSellerByPhone(phone string) (*models.Seller, error) {
	user := &models.Seller{}
	if err := pdb.DB.Where("phone_number =?", phone).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindBuyerByPhone finds a user by the phone
func (pdb PostgresDb) FindBuyerByPhone(phone string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Where("phone_number =?", phone).First(buyer).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}

// TokenInBlacklist checks if token is already in the blacklist collection
func (pdb *PostgresDb) TokenInBlacklist(token *string) bool {
	tok := &models.Blacklist{}
	if err := pdb.DB.Where("token = ?", token).First(&tok).Error; err != nil {
		return false
	}

	return true
}

// FindAllUsersExcept returns all the users expcept the one specified in the except parameter
func (pdb *PostgresDb) FindAllSellersExcept(except string) ([]models.Seller, error) {
	sellers := []models.Seller{}
	if err := pdb.DB.Not("username = ?", except).Find(sellers).Error; err != nil {

		return nil, err
	}
	return sellers, nil
}

func (pdb *PostgresDb) UpdateBuyerProfile(id uint, update *models.UpdateUser) error {
	result :=
		pdb.DB.Model(models.Buyer{}).
			Where("id = ?", id).
			Updates(
				models.User{
					FirstName:   update.FirstName,
					LastName:    update.LastName,
					PhoneNumber: update.PhoneNumber,
					Address:     update.Address,
					Email:       update.Email,
				},
			)
	return result.Error
}

func (pdb *PostgresDb) UpdateSellerProfile(id uint, update *models.UpdateUser) error {
	result :=
		pdb.DB.Model(models.Seller{}).
			Where("id = ?", id).
			Updates(
				models.User{
					FirstName:   update.FirstName,
					LastName:    update.LastName,
					PhoneNumber: update.PhoneNumber,
					Address:     update.Address,
					Email:       update.Email,
				},
			)
	return result.Error
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func (pdb *PostgresDb) UploadFileToS3(h *session.Session, file multipart.File, fileName string, size int64) (string, error) {
	// get the file size and read the file content into a buffer
	buffer := make([]byte, size)
	file.Read(buffer)
	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file you're uploading
	url := "https://s3-eu-west-3.amazonaws.com/arp-rental/" + fileName
	_, err := s3.New(h).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:                  aws.String(fileName),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	return url, err
}

func (pdb *PostgresDb) UpdateUserImageURL(username, url string) error {
	result :=
		pdb.DB.Model(models.User{}).
			Where("username = ?", username).
			Updates(
				models.User{
					Image: url,
				},
			)
	return result.Error
}
func (pdb *PostgresDb) BuyerUpdatePassword(password, newPassword string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Model(buyer).Where("password_hash =?", password).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}
func (pdb *PostgresDb) SellerUpdatePassword(password, newPassword string) (*models.Seller, error) {
	seller := &models.Seller{}
	if err := pdb.DB.Model(seller).Where("password_hash =?", password).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return seller, nil
}
func (pdb *PostgresDb) BuyerResetPassword(email, newPassword string) (*models.Buyer, error) {
	buyer := &models.Buyer{}
	if err := pdb.DB.Model(buyer).Where("email =?", email).Update("password_hash", newPassword).Error; err != nil {
		return nil, err
	}
	return buyer, nil
}

//FindIndividualSellerShop return the individual seller and its respective shop gotten by its unique ID
func (pdb *PostgresDb) FindIndividualSellerShop(sellerID string) (*models.Seller, error) {
	//create instance of a seller and its respective product, and unmarshal data into them
	seller := &models.Seller{}

	if err := pdb.DB.Preload("Product").Where("id = ?", sellerID).Find(&seller).Error; err != nil {
		log.Println("Error in finding", err)
		return nil, err
	}

	return seller, nil
}

//GetAllBuyerOrder fetches all buyer orders
func (pdb *PostgresDb) GetAllBuyerOrder(buyerId uint) ([]models.Order, error) {
	var buyerOrder []models.Order
	if err := pdb.DB.Where("buyer_id =?", buyerId).
		Preload("Seller").
		Preload("Buyer").
		Preload("Product").
		Preload("Product.Category").
		Find(&buyerOrder).
		Error; err != nil {
		log.Println("could not find order", err)
		return nil, err
	}

	return buyerOrder, nil
}

// GetAllSellerOrder fetches all buyer orders
func (pdb *PostgresDb) GetAllSellerOrder(sellerId uint) ([]models.Order, error) {
	var sellerOrder []models.Order
	if err := pdb.DB.Where("seller_id= ?", sellerId).Preload("Seller").
		Preload("Buyer").
		Preload("Product").
		Preload("Product.Category").
		Find(&sellerOrder).
		Error; err != nil {
		return nil, err
	}
	return sellerOrder, nil
}

// GetAllSellerOrderCount fetches all buyer orders
func (pdb *PostgresDb) GetAllSellerOrderCount(sellerId uint) (int, error) {
	var sellerOrder []models.Order
	if err := pdb.DB.Where("seller_id= ?", sellerId).Preload("Seller").
		Preload("Buyer").
		Preload("Product").
		Preload("Product.Category").
		Find(&sellerOrder).
		Error; err != nil {
		return 0, err
	}
	count := len(sellerOrder)

	return count, nil
}

// GetAllSellers returns all the sellers in the updated database
func (pdb *PostgresDb) GetAllSellers() ([]models.Seller, error) {
	var seller []models.Seller
	err := pdb.DB.Model(&models.Seller{}).Find(&seller).Error
	if err != nil {
		return nil, err
	}
	return seller, nil
}

// GetProductByID returns a particular product by it's ID
func (pdb *PostgresDb) GetProductByID(id uint) (*models.Product, error) {
	product := &models.Product{}
	if err := pdb.DB.Where("ID=?", id).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

//GET INDIVIDUAL SELLER PRODUCT
func (pdb *PostgresDb) FindSellerProduct(sellerID string) ([]models.Product, error) {

	product := []models.Product{}

	if err := pdb.DB.Preload("Category").Where("seller_id = ?", sellerID).Find(&product).Error; err != nil {
		log.Println("Error finding seller product", err)
		return nil, err
	}
	return product, nil

}

//GET PAID PRODUCTS FROM DATABASE
func (pdb *PostgresDb) FindPaidProduct(sellerID string) ([]models.CartProduct, error) {

	cartProduct := []models.CartProduct{}

	if err := pdb.DB.Where("order_status = ?", true).Where("seller_id = ?", sellerID).Find(&cartProduct).Error; err != nil {
		log.Println("Error finding products paid", err)
		return nil, err
	}

	return cartProduct, nil

}

func (pdb *PostgresDb) AddToCart(product models.Product, buyer *models.Buyer) error {
	var prod *models.Product
	var userBuyer *models.Buyer
	var cart *models.Cart

	err := pdb.DB.Where("id = ?", product.ID).First(&prod).Error
	if err != nil {
		return err
	}

	err = pdb.DB.Where("id = ?", buyer.ID).First(&userBuyer).Error
	if err != nil {
		return err
	}

	err = pdb.DB.Where("buyer_id = ?", buyer.ID).First(&cart).Error
	if err != nil {
		return err
	}

	cartProduct := models.CartProduct{
		CartID:        cart.ID,
		ProductID:     product.ID,
		TotalPrice:    prod.Price * product.Quantity,
		TotalQuantity: product.Quantity,
		OrderStatus:   false,
	}

	cart.Product = append(cart.Product, cartProduct)

	err = pdb.DB.Where("id = ?", cart.ID).Save(&cart).Error
	if err != nil {
		return err
	}

	return nil

}

func (pdb *PostgresDb) GetCartProducts(buyer *models.Buyer) ([]models.CartProduct, error) {

	var cart *models.Cart
	var addedProducts []models.CartProduct

	err := pdb.DB.Where("buyer_id = ?", buyer.ID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	err = pdb.DB.Where("cart_id = ?", cart.ID).Where("order_status = ?", false).
		Find(&addedProducts).Error
	if err != nil {
		return nil, err
	}

	return addedProducts, nil

}

func (pdb *PostgresDb) ViewCartProducts(addedProducts []models.CartProduct) ([]models.ProductDetails, error) {
	var details []models.ProductDetails

	for i := 0; i < len(addedProducts); i++ {
		var product *models.Product
		err := pdb.DB.Where("id = ?", addedProducts[i].ProductID).First(&product).Error
		if err != nil {
			return nil, err
		}
		prodDetail := models.ProductDetails{
			Name:     product.Title,
			Price:    addedProducts[i].TotalPrice,
			Quantity: addedProducts[i].TotalQuantity,
			Images:   product.Images,
		}

		details = append(details, prodDetail)
	}

	return details, nil
}

func (pdb *PostgresDb) DeletePaidFromCart(cartID uint) error {
	var cartProducts []models.CartProduct
	var product []models.Product

	err := pdb.DB.Where("cart_id = ?", cartID).Where("order_status = ?", false).
		Find(&cartProducts).Error
	if err != nil {
		return err
	}
	err = pdb.DB.Find(&product).Error
	if err != nil {
		return err
	}

	for i := 0; i < len(cartProducts); i++ {
		for j := i; j < len(product); j++ {
			if cartProducts[i].ProductID == product[j].ID {
				var cart models.Cart
				err := pdb.DB.Where("id = ?", cartProducts[i].CartID).First(&cart).Error
				if err != nil {
					return nil
				}
				orders := models.Order{
					SellerId:  product[j].SellerId,
					BuyerId:   cart.BuyerID,
					ProductId: product[j].ID,
				}

				newQuantity := product[j].Quantity - cartProducts[i].TotalQuantity
				err = pdb.DB.Model(&models.Product{}).Where("id=?", product[j].ID).
					Update("quantity", newQuantity).Error
				if err != nil {
					return err
				}

				err = pdb.DB.Model(&models.CartProduct{}).Where("id = ?", cartProducts[i].ID).
					Update("order_status", true).Error
				if err != nil {
					return err
				}

				err = pdb.DB.Create(&orders).Error
				if err != nil {
					return err
				}

				err = pdb.DB.Where("id = ?", cartProducts[i].ID).Delete(&models.CartProduct{}).Error
				if err != nil {
					return err
				}

				break
			}
		}
	}
	return nil
}

func (pdb *PostgresDb) GetSellersProducts(sellerID uint) ([]models.Product, error) {
	var products []models.Product

	err := pdb.DB.Where("seller_id = ?", sellerID).Find(&products).Error
	if err != nil {
		log.Println("Error from GetSellersProduct in DB")
		return nil, err
	}
	return products, nil
}
