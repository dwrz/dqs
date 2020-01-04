package diet

import "github.com/dwrz/dqs/pkg/dqs/category"

type Diet string

const (
	OMNIVORE   Diet = "omnivore"
	VEGAN      Diet = "vegan"
	VEGETARIAN Diet = "vegetarian"
)

func (d Diet) GetEntryTemplate() map[string]category.Category {
	switch d {
	case OMNIVORE:
		return omnivore
	case VEGAN:
		return vegan
	case VEGETARIAN:
		return vegetarian
	default:
		return nil
	}
}
