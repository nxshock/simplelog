package simplelog

import "time"

type UntilCalculator struct {
	calculator ProgressItf
}

func newUntilCalculator(calculator ProgressItf) *UntilCalculator {
	return &UntilCalculator{calculator}
}

// Average returns remaining time based on average speed
func (u *UntilCalculator) Average() time.Time {
	u.calculator.RLock()
	defer u.calculator.RUnlock()

	if u.calculator.Estimated() <= 0 {
		return time.Now()
	}

	return time.Now().Add(time.Second * time.Duration(float64(u.calculator.Estimated()-u.calculator.Elapsed())/u.calculator.Speed().average()))
}

// Average returns remaining time based on last complete speed
func (u *UntilCalculator) LastPeriod() time.Time {
	u.calculator.RLock()
	defer u.calculator.RUnlock()

	if u.calculator.Estimated() <= 0 || u.calculator.Elapsed() >= u.calculator.Estimated() {
		return time.Now()
	}

	result := time.Now()

	if i, ok := u.calculator.(*Progress); ok {
		return i.until
	} else if i, ok := u.calculator.(Progresses); ok {
		var maxUntilTime time.Time
		for k := range i {
			if k.until.After(maxUntilTime) {
				maxUntilTime = k.until
			}
		}

		return maxUntilTime
	}

	return result
}
