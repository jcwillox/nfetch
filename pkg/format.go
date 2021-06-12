package pkg

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

	for val, name := range map[int64]string{
		days:    "day",
		hours:   "hour",
		minutes: "minute",
	} {
		if val == 1 {
			parts = append(parts, fmt.Sprintf("%d %s", val, name))
		} else if val > 1 {
			parts = append(parts, fmt.Sprintf("%d %ss", val, name))
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
