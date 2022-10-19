package model

import "time"

type CurrencyRate struct {
	Id        int64
	Code      string
	Value     int64
	CreatedAt time.Time
}
