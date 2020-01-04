package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dwrz/dqs/pkg/dqs/category"
	"github.com/dwrz/dqs/pkg/dqs/entry"
	"github.com/dwrz/dqs/pkg/dqs/store"
)

func updatePortions(
	change string, args []string, store *store.Store, entry *entry.Entry,
) error {
	if len(args) < 2 {
		return fmt.Errorf("missing category and quantity")
	}
	if len(args)%2 != 0 {
		return fmt.Errorf("uneven number of arguments")
	}

	for i := 0; i < len(args); i += 2 {
		portionCategory := args[i]
		quantity := args[i+1]

		c, err := getCategory(entry, portionCategory)
		if err != nil {
			return err
		}

		q, err := parseQuantity(quantity)
		if err != nil {
			return err
		}

		switch change {
		case "add":
			if err := c.AddPortions(q); err != nil {
				return fmt.Errorf(
					"failed to add portions: %v", err,
				)
			}
		case "remove":
			if err := c.RemovePortions(q); err != nil {
				return fmt.Errorf(
					"failed to remove portions: %v", err,
				)
			}
		}
	}

	if err := store.UpdateEntry(entry); err != nil {
		return fmt.Errorf("failed to update store: %v", err)
	}

	return nil
}

func getCategory(entry *entry.Entry, portionCategory string) (
	*category.Category, error,
) {
	// Try expanding an abbreviation.
	c, exists := entry.Categories[category.Abbreviations[portionCategory]]
	if !exists {
		// Check if a lowercase full category name was used.
		c, exists = entry.Categories[strings.Title(portionCategory)]
		if !exists {
			// Check the full, capitalized name.
			c, exists = entry.Categories[portionCategory]
			if !exists {
				return nil, fmt.Errorf(
					"category %s not found",
					portionCategory,
				)
			}
		}
	}

	return &c, nil
}

func parseQuantity(quantity string) (float64, error) {
	q, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse quantity: %v", err)
	}
	if q < 0 {
		return 0, fmt.Errorf("cannot add negative quantity: %f", q)
	}
	if math.Mod(q, 0.5) != 0 {
		return 0, fmt.Errorf(
			"quantity %f is not a full or half portion", q,
		)
	}

	return q, nil
}
