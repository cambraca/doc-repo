package models

type User struct {
	Model
	AccountID uint     `json:"-"`
	Account   *Account `json:"account,omitempty"`
	Name      string   `gorm:"size:255" json:"name"`
	Email     string   `gorm:"size:255" json:"email"`
}
