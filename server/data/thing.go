package data

import "time"

type (
	Thing interface {
		ID() int64
		CreatedAt() time.Time
	}

	thing struct {
		id int64
		createdAt time.Time
	}
)

func (t *thing) ID() int64 {
	return t.id
}

func (t *thing) CreatedAt() time.Time {
	return t.createdAt
}