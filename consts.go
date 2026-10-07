package simplelog

import (
	"time"

	"github.com/fatih/color"
)

const (
	defaultFileTimestampFormat     = "2006-01-02 15:04:05"
	defaultTerminalTimestampFormat = "15:04:05"
	defaulLogLevel                 = LogLevelInfo
	defaultTrimMarker              = "..."

	DEFAULT_PERIOD_DURATION = time.Minute
)

var (
	defaultTimestampStyle = color.New(color.FgHiBlack).SprintFunc()
	defaultTraceStyle     = color.New(color.FgWhite).SprintFunc()
	defaultDebugStyle     = color.New(color.FgWhite).SprintFunc()
	// defaultInfoStyle uses default terminal foreground color
	defaultWarningStyle  = color.New(color.FgYellow).SprintFunc()
	defaultErrorStyle    = color.New(color.FgRed).SprintFunc()
	defaultFatalStyle    = color.New(color.FgRed).SprintFunc()
	defaultProgressStyle = color.New(color.FgHiBlack).SprintFunc()

	defaultPrintProgressFunc = func(l *Logger) {
		if tc := l.Progresses.ActiveTaskCount(); tc > 1 {
			l.Progressf("[%2d%%] %d / %d (%d tasks), %s remaining",
				l.Progresses.Percent(),
				l.Progresses.Elapsed(),
				l.Progresses.Estimated(),
				tc, l.Progresses.Remaining().LastPeriod())
		} else {
			l.Progressf("[%2d%%] %d / %d, %s remaining",
				l.Progresses.Percent(),
				l.Progresses.Elapsed(),
				l.Progresses.Estimated(),
				l.Progresses.Remaining().LastPeriod())
		}
	}
)
