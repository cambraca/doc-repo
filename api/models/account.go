package models

// An Account contains users and documents.
type Account struct {
	Model
	Name string `gorm:"size:255" json:"name"`
}
