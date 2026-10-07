package simplelog

import (
	"time"
)

type ProgressItem struct {
	logger *Logger

	periods []*Period

	periodSize time.Duration

	elapsed   uint
	estimated uint
	until     time.Time
}

// NewTimeCalc returns new time calculator with `estimated` tasks count.
func (l *Logger) NewProgress(estimated uint) *ProgressItem {
	return l.NewProgressCustom(estimated, DEFAULT_PERIOD_DURATION)
}

// NewTimeCalcCustom returns new time calculator with `estimated` tasks count and `periodSize` size of periods.
func (l *Logger) NewProgressCustom(estimated uint, periodSize time.Duration) *ProgressItem {
	now := time.Now()

	t := &ProgressItem{
		logger:     l,
		estimated:  estimated,
		periodSize: periodSize,
		periods: []*Period{{
			startTime: now,
			endTime:   now,
			processed: 0}}}

	return t
}

func (p *ProgressItem) lastCompletePeriod() *Period {
	if len(p.periods) > 1 {
		return p.periods[len(p.periods)-2]
	}

	return p.periods[len(p.periods)-1]
}

func (p *ProgressItem) tick(t time.Time, change uint) {
	p.elapsed += change

	lastPeriod := p.periods[len(p.periods)-1]
	lastPeriod.processed += change
	lastPeriod.endTime = t

	if t.Sub(lastPeriod.startTime) > p.periodSize {
		p.periods = append(p.periods, &Period{startTime: t, endTime: t, processed: 0})
	}

	if p.elapsed < p.estimated {
		p.until = t.Add(time.Second * time.Duration(float64(p.estimated-p.elapsed)/p.lastCompletePeriod().speed()))
	} else {
		p.until = time.Now()
	}
}

// IncrementProgress adds `delta“, to the task counter.
func (p *ProgressItem) IncrementProgress(delta uint) {
	now := time.Now()

	p.logger.Progresses.mu.Lock()
	p.tick(now, delta)
	p.logger.Progresses.mu.Unlock()

	// Skip progress calculation if progress message will not be displayed
	p.logger.Progresses.mu.RLock()
	if p.logger.MinProgressUpdatePeriod > 0 && time.Since(p.logger.lastProgressUpdateTime) < p.logger.MinProgressUpdatePeriod {
		p.logger.Progresses.mu.RUnlock()
		return
	}
	p.logger.Progresses.mu.RUnlock()

	p.logger.Progresses.mu.RLock()
	p.logger.PrintProgressFunc(p.logger)
	p.logger.Progresses.mu.RUnlock()
}

func (p *ProgressItem) Percent() uint {
	if p.estimated-p.elapsed == 0 {
		return 100
	}

	return p.elapsed * 100 / p.estimated
}

func (p *ProgressItem) Finish() {
	p.logger.Progresses.mu.Lock()
	delete(p.logger.Progresses.ProgressItems, p)
	p.logger.Progresses.mu.Unlock()

	p.logger.Progresses.mu.RLock()
	p.logger.PrintProgressFunc(p.logger)
	p.logger.Progresses.mu.RUnlock()
}

func (p *ProgressItem) Periods(count uint) []*Period { return p.periods }
func (p *ProgressItem) Elapsed() uint                { return p.elapsed }
func (p *ProgressItem) Estimated() uint              { return p.estimated }
