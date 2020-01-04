package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dwrz/dqs/pkg/dqs/entry"
)

const (
	DEFAULT_DB = ""
	DB_USAGE   = "the path to the dqs database file"

	DEFAULT_DATE = "today"
	DATE_USAGE   = "the entry date to use, formatted as YYYYMMDD"
)

type config struct {
	Date time.Time
	DB   string
}

func getConfig() (config, error) {
	var (
		cfg  = config{}
		date string
	)

	flag.StringVar(&cfg.DB, "database", DEFAULT_DB, DB_USAGE)
	flag.StringVar(&cfg.DB, "db", DEFAULT_DB, DB_USAGE)
	flag.StringVar(&date, "date", DEFAULT_DATE, DATE_USAGE)
	flag.StringVar(&date, "d", DEFAULT_DATE, DATE_USAGE)

	flag.Parse()

	// If unset, set the default DB path.
	if cfg.DB == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, fmt.Errorf(
				"failed to get user home dir: %v", err,
			)
		}
		cfg.DB = fmt.Sprintf("%s/.config/dqs/", home)
	}

	// Parse and set the date.
	switch date {
	case DEFAULT_DATE:
		now := time.Now()
		cfg.Date = time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)
	default:
		var err error
		cfg.Date, err = time.Parse(entry.DateFormat, date)
		if err != nil {
			return cfg, fmt.Errorf("invalid date: %v", err)
		}
	}

	return cfg, nil
}
