package simplelog

import (
	"github.com/fatih/color"
)

const (
	defaultFileTimestampFormat     = "2006-01-02 15:04:05"
	defaultTerminalTimestampFormat = "15:04:05"
	defaulLogLevel                 = LogLevelInfo
	defaultTrimMarker              = "..."
)

var (
	defaultTimestampStyle = color.New(color.FgWhite).SprintFunc()
	defaultTraceStyle     = color.New(color.FgWhite).SprintFunc()
	defaultDebugStyle     = color.New(color.FgWhite).SprintFunc()
	// defaultInfoStyle uses default terminal foreground color
	defaultWarningStyle  = color.New(color.FgYellow).SprintFunc()
	defaultErrorStyle    = color.New(color.FgRed).SprintFunc()
	defaultFatalStyle    = color.New(color.FgRed).SprintFunc()
	defaultProgressStyle = color.New(color.FgWhite).SprintFunc()
)
