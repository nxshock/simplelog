package simplelog

import (
	"sync"
	"time"
)

type Progress struct {
	ProgressItems map[*ProgressItem]struct{}

	mu sync.RWMutex
}

func (p *Progress) Elapsed() uint {
	elapsed := uint(0)

	for k := range p.ProgressItems {
		elapsed += k.elapsed
	}

	return elapsed
}

func (p *Progress) Estimated() uint {
	estimated := uint(0)

	for k := range p.ProgressItems {
		estimated += k.estimated
	}

	return estimated
}

func (p *Progress) Periods(count uint) []*Period {
	var periods []*Period

	for k := range p.ProgressItems {
		periods = append(periods, k.Periods(count)...)
	}

	return mergePeriods(periods, time.Minute) // TODO: change period size
}

func (p *Progress) ActiveTaskCount() uint {
	return uint(len(p.ProgressItems))
}

func (p *Progress) Speed() *SpeedCalculator {
	return newSpeedCalculator(p)
}

func (p *Progress) Remaining() *RemainingCalculator {
	return newRemainingCalculator(p)
}

func (p *Progress) Until() *UntilCalculator {
	return newUntilCalculator(p)
}

func (p *Progress) Percent() uint {
	if p.Estimated() <= 0 {
		return 100
	}

	return p.Elapsed() * 100 / p.Estimated()
}
