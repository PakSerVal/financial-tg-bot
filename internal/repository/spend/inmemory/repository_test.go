package inmemory

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

func Test_Save(t *testing.T) {
	repo := inmemory{}
	res, err := repo.Save(12345.0, "category")

	assert.NoError(t, err)

	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, float64(12345), res.Price)
	assert.Equal(t, "category", res.Category)
}

func Test_GetByTimeSince(t *testing.T) {
	now := time.Now()
	repo := inmemory{}

	recordHourAgo := spend.SpendRecord{
		ID:       1,
		Price:    1234.123,
		Category: "cat1",
		DateTime: now.Add(-1 * time.Hour),
	}
	recordTwoHourAgo := spend.SpendRecord{
		ID:       2,
		Price:    5678.32,
		Category: "cat2",
		DateTime: now.Add(-2 * time.Hour),
	}
	recordTwoYearsAgo := spend.SpendRecord{
		ID:       3,
		Price:    5678,
		Category: "cat3",
		DateTime: now.AddDate(-2, 0, 0),
	}
	recordMonthAgo := spend.SpendRecord{
		ID:       4,
		Price:    5678,
		Category: "cat4",
		DateTime: now.AddDate(0, -1, 0),
	}
	recordFiveDaysAgo := spend.SpendRecord{
		ID:       5,
		Price:    5678,
		Category: "cat5",
		DateTime: now.AddDate(0, 0, -5),
	}

	repo.records = []spend.SpendRecord{recordHourAgo, recordTwoHourAgo, recordTwoYearsAgo, recordMonthAgo, recordFiveDaysAgo}
	cases := []struct {
		timeSince time.Time
		wanted    []spend.SpendRecord
	}{
		{
			timeSince: time.Now().AddDate(-1, 0, 0),
			wanted:    []spend.SpendRecord{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, 0, -1),
			wanted:    []spend.SpendRecord{recordHourAgo, recordTwoHourAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -1, 0),
			wanted:    []spend.SpendRecord{recordHourAgo, recordTwoHourAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -3, 0),
			wanted:    []spend.SpendRecord{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := repo.GetByTimeSince(c.timeSince)

			assert.NoError(t, err)
			assert.Equal(t, c.wanted, res)
		})
	}
}
