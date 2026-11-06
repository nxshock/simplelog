package simplelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortify(t *testing.T) {
	var tests = []struct {
		n        uint
		expected string
	}{
		{0, "0"},
		{1000, "1.0k"},
		{1000000, "1.0M"},
		{1000000000, "1.0G"},
		{1000000000000, "1.0T"},
		{1000000000000000, "1.0P"},
		{1000000000000000000, "1.0E"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, shortify(test.n))
	}
}
