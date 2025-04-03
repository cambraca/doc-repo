package models

type Document struct {
	Model
	AccountID   int      `json:"accountId,omitempty"`
	Account     *Account `json:"account,omitempty"`
	CreatedByID int      `json:"-"`
	CreatedBy   *User    `json:"createdBy,omitempty"`
	Name        string   `gorm:"size:255" json:"name"`
	MimeType    string   `gorm:"size:40" json:"mimeType"`
	ComplexInfo struct {
		A string
		B int
	} `gorm:"serializer:json" json:"complexInfo,omitempty"`
}
