package spend

import "time"

type SpendRecord struct {
	ID       int64
	Price    float64
	Category string
	DateTime time.Time
}
