package simplelog

import "github.com/charmbracelet/lipgloss"

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelProgress LogLevel = 9
)

const (
	defaultFileTimestampFormat     = "2006-01-02 15:04:05"
	defaultTerminalTimestampFormat = "15:04:05"
	defaulLogLevel                 = LogLevelInfo
	defaultTrimMarker              = "..."
)

var (
	defaultTimestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultTraceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultDebugStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	// defaultInfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	defaultWarningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff80"))
	defaultErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	defaultFatalStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	defaultProgressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)
