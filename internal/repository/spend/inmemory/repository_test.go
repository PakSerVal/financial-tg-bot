package inmemory

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

func Test_Save(t *testing.T) {
	repo := inmemory{}
	res, err := repo.Save(1234500, "category")

	assert.NoError(t, err)

	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, int64(1234500), res.Price)
	assert.Equal(t, "category", res.Category)
}

func Test_GetByTimeSince(t *testing.T) {
	now := time.Now()
	repo := inmemory{}

	recordHourAgo := model.Spend{
		ID:       1,
		Price:    123412,
		Category: "cat1",
		DateTime: now.Add(-1 * time.Hour),
	}
	recordTwoHourAgo := model.Spend{
		ID:       2,
		Price:    567832,
		Category: "cat2",
		DateTime: now.Add(-2 * time.Hour),
	}
	recordTwoYearsAgo := model.Spend{
		ID:       3,
		Price:    567800,
		Category: "cat3",
		DateTime: now.AddDate(-2, 0, 0),
	}
	recordMonthAgo := model.Spend{
		ID:       4,
		Price:    567800,
		Category: "cat4",
		DateTime: now.AddDate(0, -1, 0),
	}
	recordFiveDaysAgo := model.Spend{
		ID:       5,
		Price:    567800,
		Category: "cat5",
		DateTime: now.AddDate(0, 0, -5),
	}

	repo.records = []model.Spend{recordHourAgo, recordTwoHourAgo, recordTwoYearsAgo, recordMonthAgo, recordFiveDaysAgo}
	cases := []struct {
		timeSince time.Time
		wanted    []model.Spend
	}{
		{
			timeSince: time.Now().AddDate(-1, 0, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, 0, -1),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -1, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordFiveDaysAgo},
		},
		{
			timeSince: time.Now().AddDate(0, -3, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
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
