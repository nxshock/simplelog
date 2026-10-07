package simplelog

import (
	"time"
)

type SpeedCalculator struct {
	calculator *Progress
}

func newSpeedCalculator(calculator *Progress) *SpeedCalculator {
	return &SpeedCalculator{calculator}
}

// Average returns average speed
func (s *SpeedCalculator) Average() float64 {
	return float64(s.calculator.Elapsed()) / time.Since(s.calculator.Periods(0)[0].startTime).Seconds()
}

// LastPeriod returns last complete period average speed
func (s *SpeedCalculator) LastPeriod() float64 {
	var lastCompletePeriod *Period
	periods := s.calculator.Periods(0)
	if len(periods) > 1 {
		lastCompletePeriod = periods[len(periods)-2]
	}

	lastCompletePeriod = periods[len(periods)-1]

	return lastCompletePeriod.speed()
}
