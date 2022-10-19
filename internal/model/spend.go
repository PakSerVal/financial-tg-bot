package model

import "time"

type Spend struct {
	Id        int64
	Price     int64
	Category  string
	UserId    int64
	CreatedAt time.Time
}
