package inmemory

import (
	"context"
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/paksergey94/telegram-bot/internal/model"
)

func Test_Save(t *testing.T) {
	repo := inmemory{
		records: map[int64][]model.Spend{},
	}
	err := repo.SaveTx(&sql.Tx{}, context.TODO(), 1234500, "category", 123)

	res := repo.records[123][0]

	assert.NoError(t, err)
	assert.Equal(t, int64(1234500), res.Price)
	assert.Equal(t, "category", res.Category)
	assert.Equal(t, int64(123), res.UserId)
}

func Test_GetByTimeSince(t *testing.T) {
	now := time.Now()
	repo := inmemory{records: map[int64][]model.Spend{}}

	recordHourAgo := model.Spend{
		Id:        1,
		Price:     123412,
		Category:  "cat1",
		CreatedAt: now.Add(-1 * time.Hour),
	}
	recordTwoHourAgo := model.Spend{
		Id:        2,
		Price:     567832,
		Category:  "cat2",
		CreatedAt: now.Add(-2 * time.Hour),
	}
	recordTwoYearsAgo := model.Spend{
		Id:        3,
		Price:     567800,
		Category:  "cat3",
		CreatedAt: now.AddDate(-2, 0, 0),
	}
	recordMonthAgo := model.Spend{
		Id:        4,
		Price:     567800,
		Category:  "cat4",
		CreatedAt: now.AddDate(0, -1, 0),
	}
	recordFiveDaysAgo := model.Spend{
		Id:        5,
		Price:     567800,
		Category:  "cat5",
		CreatedAt: now.AddDate(0, 0, -5),
	}

	repo.records = map[int64][]model.Spend{
		123: {
			recordHourAgo, recordTwoHourAgo, recordTwoYearsAgo, recordMonthAgo, recordFiveDaysAgo,
		},
	}
	cases := []struct {
		timeSince time.Time
		wanted    []model.Spend
	}{
		{
			timeSince: now.AddDate(-1, 0, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
		{
			timeSince: now.AddDate(0, 0, -1),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo},
		},
		{
			timeSince: now.AddDate(0, -1, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordFiveDaysAgo},
		},
		{
			timeSince: now.AddDate(0, -3, 0),
			wanted:    []model.Spend{recordHourAgo, recordTwoHourAgo, recordMonthAgo, recordFiveDaysAgo},
		},
	}

	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := repo.GetByTimeSince(context.TODO(), 123, c.timeSince)

			assert.NoError(t, err)
			assert.Equal(t, c.wanted, res)
		})
	}
}
