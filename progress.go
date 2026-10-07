package simplelog

import (
	"sync"
	"time"
)

type Progress struct {
	logger *Logger

	periods []*Period

	periodSize time.Duration

	elapsed   uint
	estimated uint
	until     time.Time

	// Speed     *SpeedCalculator
	// Remaining *RemainingCalculator
	// Until *UntilCalculator

	mu *sync.RWMutex
}

type ProgressItf interface {
	Periods(count uint) []*Period
	Elapsed() uint
	Estimated() uint
	Percent() uint
	Speed() *SpeedCalculator
	Remaining() *RemainingCalculator
	Until() *UntilCalculator
	RLock()
	RUnlock()
}

// NewTimeCalc returns new time calculator with `estimated` tasks count.
func (l *Logger) NewProgress(estimated uint) *Progress {
	return l.NewProgressCustom(estimated, DEFAULT_PERIOD_DURATION)
}

// NewTimeCalcCustom returns new time calculator with `estimated` tasks count and `periodSize` size of periods.
func (l *Logger) NewProgressCustom(estimated uint, periodSize time.Duration) *Progress {
	now := time.Now()

	t := &Progress{
		logger:     l,
		estimated:  estimated,
		periodSize: periodSize,
		periods: []*Period{{
			startTime: now,
			endTime:   now,
			processed: 0}},
		mu: new(sync.RWMutex)}

	// t.Speed = newSpeedCalculator(t)
	// t.Remaining = newRemainingCalculator(t)
	// t.Until = newUntilCalculator(t)

	return t
}

func (p *Progress) lastCompletePeriod() *Period {
	if len(p.periods) > 1 {
		return p.periods[len(p.periods)-2]
	}

	return p.periods[len(p.periods)-1]
}

func (p *Progress) tick(t time.Time, change uint) {
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
func (p *Progress) IncrementProgress(delta uint) {
	now := time.Now()

	p.mu.Lock()
	defer p.mu.Unlock()

	p.tick(now, delta)

	p.logger.PrintProgressFunc(p.logger)
}

func (p *Progress) RLock() {
	p.mu.RLock()
}
func (p *Progress) RUnlock() {
	p.mu.RUnlock()
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
	if p.estimated-p.elapsed == 0 {
		return 100
	}

	return p.elapsed * 100 / p.estimated
}

func (p *Progress) Finish() {
	delete(p.logger.Progresses, p)

	p.logger.PrintProgressFunc(p.logger)
}

func (p *Progress) Periods(count uint) []*Period { return p.periods }
func (p *Progress) Elapsed() uint                { return p.elapsed }
func (p *Progress) Estimated() uint              { return p.estimated }
