package simplelog

import "fmt"

var suffixes = []string{"k", "M", "G", "T", "P", "E", "Z", "Y", "R", "Q"}

func shortify(n uint) string {
	if n < 1000 {
		return fmt.Sprint(n)
	}

	f := float64(n)
	pos := -1

	for f >= 1000 && pos < len(suffixes) {
		pos += 1
		f /= 1000
	}

	return fmt.Sprintf("%.1f%s", f, suffixes[pos])
}
