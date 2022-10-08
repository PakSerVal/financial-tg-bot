package spend

import "time"

type SpendRecord struct {
	ID       int64
	Price    int64
	Category string
	DateTime time.Time
}
