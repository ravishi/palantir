package data

import "time"

type (
	Thing struct {
		ID        int64 `gorm:"primary_key" json:"id"`
		CreatedAt time.Time
	}
)
