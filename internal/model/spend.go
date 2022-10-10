package model

import "time"

type Spend struct {
	ID       int64
	Price    float64
	Category string
	DateTime time.Time
}
