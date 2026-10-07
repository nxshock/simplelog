package simplelog

import "time"

type UntilCalculator struct {
	calculator *Progress
}

func newUntilCalculator(calculator *Progress) *UntilCalculator {
	return &UntilCalculator{calculator}
}

// Average returns remaining time based on average speed
func (u *UntilCalculator) Average() time.Time {
	if u.calculator.Estimated() <= 0 {
		return time.Now()
	}

	return time.Now().Add(time.Second * time.Duration(float64(u.calculator.Estimated()-u.calculator.Elapsed())/u.calculator.Speed().Average()))
}

// Average returns remaining time based on last complete speed
func (u *UntilCalculator) LastPeriod() time.Time {
	if u.calculator.Estimated() <= 0 || u.calculator.Elapsed() >= u.calculator.Estimated() {
		return time.Now()
	}

	result := time.Now()

	for k := range u.calculator.ProgressItems {
		if k.until.After(result) {
			result = k.until
		}
	}

	return result
}
