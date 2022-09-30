package inmemory

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/repository/spend"
)

func Test_Save(t *testing.T) {
	repo := New()
	res, err := repo.Save(12345, "category")

	assert.NoError(t, err)

	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, int64(12345), res.Sum)
	assert.Equal(t, "category", res.Category)
}

func Test_GetByTimeSince(t *testing.T) {
	now := time.Now()
	repo := New()

	recordHourAgo := spend.Record{
		ID:       1,
		Sum:      1234,
		Category: "cat1",
		DateTime: now.Add(-1 * time.Hour),
	}
	recordTwoHourAgo := spend.Record{
		ID:       2,
		Sum:      5678,
		Category: "cat2",
		DateTime: now.Add(-2 * time.Hour),
	}
	recordTwoYearsAgo := spend.Record{
		ID:       3,
		Sum:      5678,
		Category: "cat3",
		DateTime: now.AddDate(-2, 0, 0),
	}
	recordMonthAgo := spend.Record{
		ID:       4,
		Sum:      5678,
		Category: "cat4",
		DateTime: now.AddDate(0, -1, 0),
	}
	recordFiveDaysAgo := spend.Record{
		ID:       5,
		Sum:      5678,
		Category: "cat5",
		DateTime: now.AddDate(0, 0, -5),
	}

	repo.records = []spend.Record{recordHourAgo, recordTwoHourAgo, recordTwoYearsAgo, recordMonthAgo, recordFiveDaysAgo}
	cases := []struct {
		timeSince time.Time
		wanted    []spend.Record
	}{
		{
			timeSince: time.Now().AddDate(-1, 0, 0),
			wanted:    []spend.Record{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, 0, -1),
			wanted:    []spend.Record{recordHourAgo, recordTwoHourAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -1, 0),
			wanted:    []spend.Record{recordHourAgo, recordTwoHourAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -3, 0),
			wanted:    []spend.Record{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
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
