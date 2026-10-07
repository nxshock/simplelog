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

func mergePeriods(p []*Period, periodSize time.Duration) []*Period {
	resultMap := make(map[Period]*Period)

	for _, v := range p {
		tmpPeriod := Period{startTime: v.startTime.Truncate(periodSize), endTime: v.endTime.Truncate(periodSize)}
		period := resultMap[tmpPeriod]
		if period == nil {
			period = &Period{startTime: tmpPeriod.startTime, endTime: tmpPeriod.endTime}
			resultMap[tmpPeriod] = period
		}
		period.processed += v.processed
	}

	result := make([]*Period, 0, len(resultMap))
	for _, v := range resultMap {
		result = append(result, v)
	}

	return result
}
