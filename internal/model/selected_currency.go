package model

import "time"

type SelectedCurrency struct {
	Id        int64
	Code      string
	UserId    int64
	CreatedAt time.Time
}
