package simplelog

import (
	"time"
)

type SpeedCalculator struct {
	calculator *TimeCalc
}

func newSpeedCalculator(calculator *TimeCalc) *SpeedCalculator {
	return &SpeedCalculator{calculator}
}

func (s *SpeedCalculator) average() float64 {
	return float64(s.calculator.elapsed) / time.Since(s.calculator.periods[0].startTime).Seconds()
}

// Average returns average speed
func (s *SpeedCalculator) Average() float64 {
	s.calculator.mu.RLock()
	defer s.calculator.mu.RUnlock()

	return s.average()
}

// LastPeriod returns last complete period average speed
func (s *SpeedCalculator) LastPeriod() float64 {
	s.calculator.mu.RLock()
	defer s.calculator.mu.RUnlock()

	return s.calculator.lastCompletePeriod().speed()
}
