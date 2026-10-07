package simplelog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMergePeriods(t *testing.T) {
	periods := []*Period{
		{
			startTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
			processed: 1},
		{
			startTime: time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),
			endTime:   time.Date(2000, 1, 1, 0, 0, 3, 0, time.UTC),
			processed: 2},
	}

	expectedPeriods := []*Period{
		{
			startTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			processed: 3},
	}

	assert.EqualValues(t, expectedPeriods, mergePeriods(periods, time.Minute))
}
