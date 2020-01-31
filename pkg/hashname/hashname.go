package hashname

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// HashToName converts a hash string (or any hex string > 6 chars) into a human-readable name
func HashToName(hash string) string {
	val, err := strconv.ParseInt(hash[:6], 16, 32)
	if err != nil {
		return "fudged-bug-error"
	}
	high := (val & 0xFFF000) >> 12
	low := val & 0xFFF

	return adjectives[low] + "-" + names[high]
}

// NameToHash converts a name from HashToName back into a six-character hex string which can be
// used to identify (f.e.) a commit hash.
func NameToHash(name string) (string, error) {
	// Adjective might have hyphens but the name doesn't. Split and merge
	parts := strings.Split(name, "-")
	if len(parts) < 2 {
		return "000000", errors.New("incorrect formatting for name")
	}
	person := parts[len(parts)-1]
	adj := strings.Join(parts[:len(parts)-1], "-")

	high := 0
	for i, v := range names {
		if v == person {
			high = i
			break
		}
	}
	low := 0
	for i, v := range adjectives {
		if v == adj {
			low = i
			break
		}
	}

	return fmt.Sprintf("%06x", ((high&0xFFF)<<12)|(low&0xFFF)), nil
}
