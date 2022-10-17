package model

import "time"

type Budget struct {
	Id        int64
	UserId    int64
	Value     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
