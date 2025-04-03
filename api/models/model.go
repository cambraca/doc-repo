package models

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	//DeletedAt DeletedAt `gorm:"index"`
}

//func (m *Model) MarshalJSON() ([]byte, error) {
//	data := map[string]interface{}{
//		"id":         m.ID,
//		"type":       "document", // TODO
//		"attributes": *m,
//	}
//
//	bytes, err := json.Marshal(m)
//	if err != nil {
//		panic("Could not marshal JSON for model")
//	}
//
//	ret := map[string]interface{}{
//		"a": 3,
//	}
//	return nil, nil
//}
