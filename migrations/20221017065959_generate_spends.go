package migrations

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
)

const userId = 1099595594

var categories = []string{"такси", "продукты", "техника", "инвестиции", "развлечение", "развитие"}

func init() {
	goose.AddMigration(upGenerateSpends, downGenerateSpends)
}

func upGenerateSpends(tx *sql.Tx) error {
	query := `insert into spend (price, category, user_id, created_at) values ($1, $2, $3, $4)`

	g := errgroup.Group{}
	for i := 0; i < 1000; i++ {
		g.Go(func() error {
			rand.Seed(time.Now().UnixNano())
			_, err := tx.Exec(query, rand.Intn(300000)+100, categories[rand.Intn(len(categories))], userId, generateRandomDate())

			return err
		})
	}

	return g.Wait()
}

func downGenerateSpends(tx *sql.Tx) error {
	_, err := tx.Exec(`truncate table spend`)
	return err
}

func generateRandomDate() time.Time {
	min := time.Date(2021, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Add(-time.Hour * 24 * 31).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
