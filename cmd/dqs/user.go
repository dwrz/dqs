package main

import (
	"fmt"

	"github.com/dwrz/dqs/pkg/dqs/diet"
	"github.com/dwrz/dqs/pkg/dqs/store"
	"github.com/dwrz/dqs/pkg/dqs/user"
)

func updateUser(args []string, store *store.Store, user *user.User) error {
	if len(args) == 0 {
		return fmt.Errorf("missing subcommand")
	}

	subcommand := args[0]

	switch subcommand {
	case "diet":
		if len(args) < 2 {
			return fmt.Errorf("missing diet")
		}
		diet := args[1]

		if err := setDiet(diet, store, user); err != nil {
			return fmt.Errorf("failed to set diet: %v", err)
		}
	default:
		return fmt.Errorf("unrecognized subcommand: %s", subcommand)
	}

	return nil
}

func setDiet(d string, store *store.Store, user *user.User) error {
	switch diet.Diet(d) {
	case diet.OMNIVORE:
		user.Diet = diet.OMNIVORE
	case diet.VEGAN:
		user.Diet = diet.VEGAN
	case diet.VEGETARIAN:
		user.Diet = diet.VEGETARIAN
	default:
		return fmt.Errorf("unrecognized diet: %s", d)
	}

	if err := store.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update store: %v", err)
	}

	return nil
}
