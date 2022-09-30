package spend

import "time"

type Record struct {
	ID       int64
	Sum      int64
	Category string
	DateTime time.Time
}
