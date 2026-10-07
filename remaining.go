package simplelog

import "time"

type RemainingCalculator struct {
	calculator ProgressItf
}

func newRemainingCalculator(calculator ProgressItf) *RemainingCalculator {
	return &RemainingCalculator{calculator}
}

// Average returns remaining time based on average speed
func (r *RemainingCalculator) Average() time.Duration {
	r.calculator.RLock()
	defer r.calculator.RUnlock()

	if r.calculator.Estimated() > 0 {
		return max(
			time.Duration(float64(r.calculator.Estimated()-r.calculator.Elapsed())/r.calculator.Speed().average())*time.Second,
			0)
	}

	return 0
}

// Average returns remaining time based on last complete speed
func (r *RemainingCalculator) LastPeriod() time.Duration {
	r.calculator.RLock()
	defer r.calculator.RUnlock()

	return max(time.Until(r.calculator.Until().LastPeriod()), 0)
}
