package models

type User struct {
	Model
	FirstName string `gorm:"type:varchar(25)" json:"first_name"`
	LastName  string `gorm:"type:varchar(25)" json:"last_name"`
	Email     string `gorm:"type:varchar(320);unique" json:"email"`
}
