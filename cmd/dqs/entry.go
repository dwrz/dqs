package main

import (
	"fmt"
	"time"

	"github.com/dwrz/dqs/pkg/dqs/entry"
	"github.com/dwrz/dqs/pkg/dqs/store"
	"github.com/dwrz/dqs/pkg/dqs/user"
)

func getEntry(date time.Time, store *store.Store, u *user.User) (
	*entry.Entry, error,
) {
	e, err := store.GetEntry(date.Format(entry.DateFormat))
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %v", err)
	}
	if e == nil {
		// If there's no entry for this date, return an empty template.
		e = &entry.Entry{
			Date:       date,
			Categories: u.Diet.GetEntryTemplate(),
		}
	}

	return e, nil
}

func updateEntry(args []string, store *store.Store, entry *entry.Entry) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand")
	}

	subcommand := args[0]

	switch subcommand {
	case "delete":
		if err := store.DeleteEntry(entry.GetKey()); err != nil {
			return fmt.Errorf("failed to delete entry: %v", err)
		}
	case "note":
		if len(args) < 2 {
			return fmt.Errorf("missing note")
		}
		note := args[1]

		if err := setNote(store, entry, note); err != nil {
			return fmt.Errorf("failed to set note: %v", err)
		}

	default:
		return fmt.Errorf("unrecognized subcommand: %s", subcommand)
	}

	return nil
}

func setNote(store *store.Store, entry *entry.Entry, note string) error {
	entry.Note = note

	if err := store.UpdateEntry(entry); err != nil {
		return fmt.Errorf("failed to update store: %v", err)
	}

	return nil
}
