package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func ToHumanTime(seconds uint64) string {
	duration := time.Duration(seconds) * time.Second
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))

	// show seconds only if the duration is under a minute
	if seconds < 60 {
		if seconds == 1 {
			return "1 second"
		} else {
			return fmt.Sprintf("%d seconds", seconds)
		}
	}

	parts := make([]string, 0, 3)

	for _, item := range []struct {
		Name  string
		Value int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
	} {
		if item.Value == 1 {
			parts = append(parts, fmt.Sprintf("%d %s", item.Value, item.Name))
		} else if item.Value > 1 {
			parts = append(parts, fmt.Sprintf("%d %ss", item.Value, item.Name))
		}
	}

	return strings.Join(parts, ", ")
}

func BytesToHuman(val float64, precision uint, precisionStart string) (string, string) {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	var unit string
	var usePrecision bool
	for _, unit = range units {
		if unit == precisionStart {
			usePrecision = true
		}
		if val < 1024 {
			break
		}
		val /= 1024
	}
	if usePrecision {
		return fmt.Sprintf("%.*f", precision, val), unit
	}
	return fmt.Sprintf("%.f", val), unit
}
