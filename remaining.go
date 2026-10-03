package simplelog

import "time"

type RemainingCalculator struct {
	calculator *TimeCalc
}

func newRemainingCalculator(calculator *TimeCalc) *RemainingCalculator {
	return &RemainingCalculator{calculator}
}

// Average returns remaining time based on average speed
func (r *RemainingCalculator) Average() time.Duration {
	r.calculator.mu.RLock()
	defer r.calculator.mu.RUnlock()

	if r.calculator.estimated > 0 {
		return max(
			time.Duration(float64(r.calculator.estimated-r.calculator.elapsed)/r.calculator.Speed.average())*time.Second,
			0)
	}

	return 0
}

// Average returns remaining time based on last complete speed
func (r *RemainingCalculator) LastPeriod() time.Duration {
	r.calculator.mu.RLock()
	defer r.calculator.mu.RUnlock()

	return max(time.Until(r.calculator.until), 0)
}
