package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dwrz/dqs/pkg/dqs/store"
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad config: %v\n", err)
		return
	}

	args := flag.Args()

	store, err := store.Open(cfg.DB)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer store.Close()

	// Retrieve user and entry for date.
	u, err := store.GetUser()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	e, err := getEntry(cfg.Date, store, u)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// If no arguments were provided, just print the entry.
	if len(args) == 0 {
		fmt.Println(e.FormatPrint())
		return
	}
	switch args[0] {
	case "add":
		if err := updatePortions(
			"add", args[1:], store, e,
		); err != nil {
			fmt.Fprintf(
				os.Stderr, "failed to add portions: %v\n", err,
			)
			return
		}

		fmt.Println(e.FormatPrint())

	case "entry":
		if err := updateEntry(args[1:], store, e); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

	case "remove":
		if err := updatePortions(
			"remove", args[1:], store, e,
		); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to remove portions: %v\n", err,
			)
			return
		}

		fmt.Println(e.FormatPrint())

	case "user":
		if err := updateUser(args[1:], store, u); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

	default:
		fmt.Fprintf(os.Stderr, "unrecognized command: %s\n", os.Args[1])
	}
}
