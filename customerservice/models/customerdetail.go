package models

import "time"

// UserModel represents the signup information for a user.
// Fields:
// - Username: unique username/login
// - Password: hashed password
// - UserTypeID: foreign key referencing UserTypeModel.UserTypeId
// - UserType: association to the UserTypeModel
type CustomerDetail struct {
	CustomerId   int        `json:"customer_id" gorm:"primaryKey;autoIncrement:true"`
	Name         string     `json:"name" gorm:"not null"`
	EmailAddress string     `json:"email_address" gorm:"unique;not null"`
	PhoneNumber  string     `json:"phonenumber" gorm:"not null"`
	Password     string     `json:"password" gorm:"not null"`
	CreateAt     *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
}

type UserLoginModel struct {
	EmailAddress string `json:"email_address"`
	Password     string `json:"password"`
}
