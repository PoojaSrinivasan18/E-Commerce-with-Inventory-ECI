package user

import (
	database "customerservice/database"
	models "customerservice/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/martian/log"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new customer
// @Description Create a new customer account
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CustomerDetail true "Customer registration details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /customersignup [post]
func AddNewCustomer(c *gin.Context) {
	var userSignUpModel models.CustomerDetail
	if err := c.ShouldBind(&userSignUpModel); err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if userSignUpModel.Name == "" || userSignUpModel.EmailAddress == "" || userSignUpModel.Password == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Some of the fields are not having right values"})
		return
	}

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userSignUpModel.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf("password hash error %v", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error processing password"})
		return
	}
	userSignUpModel.Password = string(hashedPassword)
	userSignUpModel.CreateAt = func(t time.Time) *time.Time { return &t }(time.Now())

	db := database.GetDB()

	var count int64
	if err := db.Model(&models.CustomerDetail{}).
		Where("email_address = ?", userSignUpModel.EmailAddress).
		Count(&count).Error; err != nil {
		log.Errorf("DB count error %v", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	if count > 0 {
		c.IndentedJSON(http.StatusConflict, gin.H{"message": "Email address already exists"})
		return
	}

	tx := db.Create(&userSignUpModel)
	if tx.Error != nil {
		log.Errorf("DB create error %v", tx.Error)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving data"})
		return
	}

	userSignUpModel.Password = "" // Do not return password
	c.IndentedJSON(http.StatusOK, "user created successfully.")
}

// @Summary Customer login
// @Description Authenticate a user and return a JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginModel true "Login credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /customerlogin [post]
func CustomerLogin(c *gin.Context) {
	var userLoginModel models.UserLoginModel
	if err := c.ShouldBind(&userLoginModel); err != nil {
		log.Errorf("FORM binding error %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var existingUser models.CustomerDetail
	db := database.GetDB()

	// Lookup user by username only (password is hashed in DB)
	if err := db.Where("email_address = ?", userLoginModel.EmailAddress).First(&existingUser).Error; err != nil {
		log.Errorf("DB query error %v", err)
		// Do not reveal whether username or password was wrong
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	// Compare provided password with hashed password stored in DB
	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userLoginModel.Password)); err != nil {
		log.Errorf("password mismatch %v", err)
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "JWT_SECRET" // replace
	}

	claims := jwt.MapClaims{
		"sub":           existingUser.CustomerId,
		"email_address": existingUser.EmailAddress,
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Errorf("token sign error %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}

	// Do not include password in response
	existingUser.Password = ""

	c.IndentedJSON(http.StatusOK, gin.H{
		"access_token": tokenString,
	})
}
