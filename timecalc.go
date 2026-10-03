package simplelog

import (
	"sync"
	"time"
)

type TimeCalc struct {
	periods []*Period

	periodSize time.Duration

	elapsed   uint
	estimated uint
	until     time.Time

	Speed     *SpeedCalculator
	Remaining *RemainingCalculator
	Until     *UntilCalculator

	mu sync.RWMutex
}

// NewTimeCalc returns new time calculator with `estimated` tasks count.
func NewTimeCalc(estimated uint) *TimeCalc {
	return NewTimeCalcCustom(estimated, DEFAULT_PERIOD_DURATION)
}

// NewTimeCalcCustom returns new time calculator with `estimated` tasks count and `periodSize` size of periods.
func NewTimeCalcCustom(estimated uint, periodSize time.Duration) *TimeCalc {
	now := time.Now()

	t := &TimeCalc{
		estimated:  estimated,
		periodSize: periodSize,
		periods: []*Period{{
			startTime: now,
			endTime:   now,
			processed: 0}}}

	t.Speed = newSpeedCalculator(t)
	t.Remaining = newRemainingCalculator(t)
	t.Until = newUntilCalculator(t)

	return t
}

func (c *TimeCalc) lastCompletePeriod() *Period {
	if len(c.periods) > 1 {
		return c.periods[len(c.periods)-2]
	}

	return c.periods[len(c.periods)-1]
}

func (c *TimeCalc) tick(t time.Time, change uint) {
	c.elapsed += change

	lastPeriod := c.periods[len(c.periods)-1]
	lastPeriod.processed += change
	lastPeriod.endTime = t

	if t.Sub(lastPeriod.startTime) > c.periodSize {
		c.periods = append(c.periods, &Period{startTime: t, endTime: t, processed: 0})
	}

	c.until = t.Add(time.Second * time.Duration(float64(c.estimated-c.elapsed)/c.lastCompletePeriod().speed()))
}

// IncrementProgress adds `delta“, to the task counter.
func (c *TimeCalc) IncrementProgress(delta uint) {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.tick(now, delta)
}
