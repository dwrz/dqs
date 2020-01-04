package entry

import (
	"fmt"

	"sort"
	"strings"
	"time"

	"github.com/dwrz/dqs/pkg/color"
	"github.com/dwrz/dqs/pkg/dqs/category"
)

const (
	DateFormat        = "20060102"
	dateDisplayFormat = "2006-01-02"
	hr                = "--------------------------------------------"
)

type Entry struct {
	Date       time.Time                    `json:"date"`
	Categories map[string]category.Category `json:"categories"`
	Note       string                       `json:"note"`
	Weight     int                          `json:"weight"`
}

// GetKey returns the key used to retrieve an entry from the store.
func (e *Entry) GetKey() string {
	return e.Date.Format(DateFormat)
}

func (e *Entry) CalculateScore() (total float64) {
	for _, c := range e.Categories {
		total += c.CalculateScore()
	}

	return total
}

// FormatPrint formats an entry for display to the user.
func (e *Entry) FormatPrint() string {
	// Assemble the data for display.
	var (
		highQuality []category.Category
		lowQuality  []category.Category
	)

	// Separate high quality and low quality categories.
	for _, c := range e.Categories {
		if c.HighQuality {
			highQuality = append(highQuality, c)
			continue
		}
		lowQuality = append(lowQuality, c)
	}

	// Sort alphabetically for consistent appearance.
	sort.Slice(highQuality, func(i, j int) bool {
		return highQuality[i].Name < highQuality[j].Name
	})
	sort.Slice(lowQuality, func(i, j int) bool {
		return lowQuality[i].Name < lowQuality[j].Name
	})

	// Prepare a string for display to the user.
	var str strings.Builder

	// Format the date.
	str.WriteString(fmt.Sprintf("|%-44s|\n", hr))
	str.WriteString(fmt.Sprintf(
		"|%-22s|\n", centerPad(e.Date.Format(dateDisplayFormat), 44),
	))
	str.WriteString(fmt.Sprintf("|%-44s|\n", hr))

	// Format high-quality categories.
	str.WriteString(fmt.Sprintf("|%-44s|\n", "High Quality"))
	str.WriteString(fmt.Sprintf("|%-44s|\n",
		hr))
	for _, c := range highQuality {
		str.WriteString(c.FormatPrint())
		str.WriteString("\n")
	}

	// Format low-quality categories.
	str.WriteString(fmt.Sprintf("|%-44s|\n", hr))
	str.WriteString(fmt.Sprintf("|%-44s|\n", "Low Quality"))
	str.WriteString(fmt.Sprintf("|%-44s|\n", hr))
	for _, c := range lowQuality {
		str.WriteString(c.FormatPrint())
		str.WriteString("\n")
	}
	str.WriteString(fmt.Sprintf("|%-44s|\n", hr))

	// Format the total.
	var totalColor color.Color
	total := e.CalculateScore()
	switch {
	case total >= 15:
		totalColor = color.BrightGreen
	case total >= 0:
		totalColor = color.BrightYellow
	default:
		totalColor = color.BrightRed

	}
	str.WriteString(fmt.Sprintf(
		"Total: %s%.1f%s\n", totalColor, total, color.Reset,
	))

	// If a note is set, format it.
	if e.Note != "" {
		str.WriteString(fmt.Sprintf("Note: %s\n", e.Note))
	}

	return str.String()
}

func centerPad(s string, width int) string {
	return fmt.Sprintf(
		"%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(s))/2, s),
	)
}
