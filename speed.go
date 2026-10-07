package simplelog

import (
	"time"
)

type SpeedCalculator struct {
	calculator ProgressItf
}

func newSpeedCalculator(calculator ProgressItf) *SpeedCalculator {
	return &SpeedCalculator{calculator}
}

func (s *SpeedCalculator) average() float64 {
	return float64(s.calculator.Elapsed()) / time.Since(s.calculator.Periods(0)[0].startTime).Seconds()
}

// Average returns average speed
func (s *SpeedCalculator) Average() float64 {
	s.calculator.RLock()
	defer s.calculator.RUnlock()

	return s.average()
}

// LastPeriod returns last complete period average speed
func (s *SpeedCalculator) LastPeriod() float64 {
	s.calculator.RLock()
	defer s.calculator.RUnlock()

	var lastCompletePeriod *Period
	periods := s.calculator.Periods(0)
	if len(periods) > 1 {
		lastCompletePeriod = periods[len(periods)-2]
	}

	lastCompletePeriod = periods[len(periods)-1]

	return lastCompletePeriod.speed()
}
