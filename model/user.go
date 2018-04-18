package model

type User struct {
	Email         string `gorm:"unique_index;primary_key"`
	FirstName     string `gorm:"not null"`
	LastName      string `gorm:"not null"`
	Password      string `gorm:"not null"`
	EthAddr       string `gorm:"unique_index;not null"`
	BtcAddr       string
	EmailVerified bool `gorm:"not null"`
	KYCVerified   bool `gorm:"not null"`
	Enable        bool `gorm:"not null"`
	IsAdmin       bool `gorm:"not null"`
}
