package simplelog

import "time"

type UntilCalculator struct {
	calculator *TimeCalc
}

func newUntilCalculator(calculator *TimeCalc) *UntilCalculator {
	return &UntilCalculator{calculator}
}

// Average returns remaining time based on average speed
func (u *UntilCalculator) Average() time.Time {
	u.calculator.mu.RLock()
	defer u.calculator.mu.RUnlock()

	if u.calculator.estimated <= 0 {
		return time.Now()
	}

	return time.Now().Add(time.Second * time.Duration(float64(u.calculator.estimated-u.calculator.elapsed)/u.calculator.Speed.average()))
}

// Average returns remaining time based on last complete speed
func (u *UntilCalculator) LastPeriod() time.Time {
	u.calculator.mu.RLock()
	defer u.calculator.mu.RUnlock()

	if u.calculator.estimated <= 0 {
		return time.Now()
	}

	return u.calculator.until
}
