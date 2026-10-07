package simplelog

import "time"

type RemainingCalculator struct {
	calculator *Progress
}

func newRemainingCalculator(calculator *Progress) *RemainingCalculator {
	return &RemainingCalculator{calculator}
}

// Average returns remaining time based on average speed
func (r *RemainingCalculator) Average() time.Duration {
	if r.calculator.Estimated() > 0 {
		return max(
			time.Duration(float64(r.calculator.Estimated()-r.calculator.Elapsed())/r.calculator.Speed().Average())*time.Second,
			0)
	}

	return 0
}

// Average returns remaining time based on last complete speed
func (r *RemainingCalculator) LastPeriod() time.Duration {
	return max(time.Until(r.calculator.Until().LastPeriod()), 0)
}
