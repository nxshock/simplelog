package simplelog

import "time"

type Period struct {
	startTime time.Time
	endTime   time.Time
	processed uint
}

func (p *Period) speed() float64 {
	periodLen := p.endTime.Sub(p.startTime).Seconds()
	if periodLen == 0 {
		return 0
	}

	return float64(p.processed) / periodLen
}
