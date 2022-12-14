package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	_ "gitlab.ozon.dev/paksergey94/telegram-bot/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	_ = flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 3 {
		flags.Usage()
		return
	}

	dbstring, command := args[1], args[2]

	db, err := goose.OpenDBWithDriver("postgres", dbstring)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
