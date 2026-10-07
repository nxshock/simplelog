package simplelog

type Progresses map[*Progress]struct{}

func (p Progresses) Elapsed() uint {
	elapsed := uint(0)

	for k := range p {
		elapsed += k.elapsed
	}

	return elapsed
}

func (p Progresses) Estimated() uint {
	estimated := uint(0)

	for k := range p {
		estimated += k.estimated
	}

	return estimated
}

func (p Progresses) Periods(count uint) []*Period {
	var periods []*Period

	for k := range p {
		periods = append(periods, k.Periods(count)...)
	}

	return periods
}

func (p Progresses) ActiveTaskCount() uint {
	return uint(len(p))
}

func (p Progresses) RLock()   {}
func (p Progresses) RUnlock() {}

func (p Progresses) Speed() *SpeedCalculator {
	return newSpeedCalculator(p)
}

func (p Progresses) Remaining() *RemainingCalculator {
	return newRemainingCalculator(p)
}

func (p Progresses) Until() *UntilCalculator {
	return newUntilCalculator(p)
}

func (p Progresses) Percent() uint {
	if p.Estimated()-p.Elapsed() == 0 {
		return 100
	}

	return p.Elapsed() * 100 / p.Estimated()

}
