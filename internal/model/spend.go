package model

import "time"

type Spend struct {
	ID       int64
	Price    int64
	Category string
	DateTime time.Time
}
