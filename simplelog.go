package simplelog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

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

const defaulLogLevel = LogLevelInfo

var (
	defaultFileTimestampFormat     = "2006-01-02 15:04:05"
	defaultTerminalTimestampFormat = "15:04:05"

	defaultTimestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultTraceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultDebugStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	// defaultInfoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	defaultWarningStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff80"))
	defaultErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	defaultFatalStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	defaultProgressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
)

type Logger struct {
	Writer io.Writer

	// timestamp format
	TimeFormat string

	// timestamp style
	TimeStampStyle lipgloss.Style

	// log level styles
	Styles map[LogLevel]*lipgloss.Style

	// strip message from spaces before output
	StripMessages bool

	// is output to terminal
	isTerminal bool

	// last written progress message length
	lastProgressLineWidth int

	Level LogLevel

	// mutex to prevent race conditions
	mu *sync.Mutex
}

type msg struct {
	TimeStamp string
	Prefix    string
	Text      string
}

func (m *msg) String() string {
	sb := new(strings.Builder)

	if m.TimeStamp != "" {
		sb.WriteString(m.TimeStamp)
		sb.WriteRune(' ')
	}

	if m.Prefix != "" {
		sb.WriteString(m.Prefix)
		sb.WriteRune(' ')
	}

	sb.WriteString(m.Text)

	return sb.String()
}

func NewLogger(w io.Writer) *Logger {
	logger := &Logger{
		Writer:         w,
		TimeStampStyle: defaultTimestampStyle,
		Styles:         make(map[LogLevel]*lipgloss.Style),
		Level:          defaulLogLevel,
		mu:             new(sync.Mutex)}

	logger.Styles[LogLevelTrace] = &defaultTraceStyle
	logger.Styles[LogLevelDebug] = &defaultDebugStyle
	// logger.Styles[LogLevelInfo] = &defaultInfoStyle
	logger.Styles[LogLevelWarn] = &defaultWarningStyle
	logger.Styles[LogLevelError] = &defaultErrorStyle
	logger.Styles[LogLevelFatal] = &defaultFatalStyle
	logger.Styles[LogLevelProgress] = &defaultProgressStyle

	if f, ok := w.(*os.File); ok {
		logger.isTerminal = term.IsTerminal(int(f.Fd()))
	}

	if logger.isTerminal {
		logger.TimeFormat = defaultTerminalTimestampFormat
	} else {
		logger.TimeFormat = defaultFileTimestampFormat
	}

	return logger
}

func minNotLessThanZero(a, b int) int {
	tmp := min(a, b)

	if tmp < 0 {
		return 0
	}

	return tmp
}

func levelSymbol(logLevel LogLevel) string {
	switch logLevel {
	case LogLevelTrace:
		return "TRC"
	case LogLevelDebug:
		return "DBG"
	case LogLevelInfo:
		return "INF"
	case LogLevelWarn:
		return "WRN"
	case LogLevelError:
		return "ERR"
	case LogLevelFatal:
		return "FTL"
	}

	return "???"
}

func (l *Logger) GetWidth() int {
	if !l.isTerminal {
		return 0
	}

	f, ok := l.Writer.(*os.File)
	if !ok {
		return 0
	}

	w, _, err := term.GetSize(int(f.Fd()))
	if err != nil {
		return 0
	}

	return w
}

func (l *Logger) Trace(a ...any) (n int, err error) {
	return l.Print(LogLevelTrace, a...)
}

func (l *Logger) Debug(a ...any) (n int, err error) {
	return l.Print(LogLevelDebug, a...)
}

func (l *Logger) Info(a ...any) (n int, err error) {
	return l.Print(LogLevelInfo, a...)
}

func (l *Logger) Warn(a ...any) (n int, err error) {
	return l.Print(LogLevelWarn, a...)
}

func (l *Logger) Error(a ...any) (n int, err error) {
	return l.Print(LogLevelError, a...)
}

func (l *Logger) Fatal(a ...any) {
	l.mu.Lock()
	l.Print(LogLevelFatal, a...)
	l.mu.Unlock()

	os.Exit(1)
}

func (l *Logger) Traceln(a ...any) (n int, err error) {
	return l.Println(LogLevelTrace, a...)
}

func (l *Logger) Debugln(a ...any) (n int, err error) {
	return l.Println(LogLevelDebug, a...)
}

func (l *Logger) Infoln(a ...any) (n int, err error) {
	return l.Println(LogLevelInfo, a...)
}

func (l *Logger) Warnln(a ...any) (n int, err error) {
	return l.Println(LogLevelWarn, a...)
}

func (l *Logger) Errorln(a ...any) (n int, err error) {
	return l.Println(LogLevelError, a...)
}

func (l *Logger) Fatalln(a ...any) {
	l.Println(LogLevelError, a...)

	if f, ok := l.Writer.(*os.File); ok {
		f.Sync()
	}

	os.Exit(1)
}

func (l *Logger) Tracef(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelInfo, format, a...)
}

func (l *Logger) Debugf(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelInfo, format, a...)
}

func (l *Logger) Infof(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelInfo, format, a...)
}

func (l *Logger) Warnf(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelWarn, format, a...)
}

func (l *Logger) Errorf(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelError, format, a...)
}

func (l *Logger) Fatalf(format string, a ...any) {
	l.Printf(LogLevelFatal, format, a...)

	os.Exit(1)
}

func (l *Logger) timestamp() string {
	if l.TimeFormat == "" {
		return ""
	}

	if !l.isTerminal {
		return time.Now().Format(l.TimeFormat)
	}

	return l.TimeStampStyle.Render(time.Now().Format(l.TimeFormat))
}

func (l *Logger) prefix(logLevel LogLevel) string {
	return fmt.Sprintf("|%s|", levelSymbol(logLevel))
}

func (l *Logger) Print(logLevel LogLevel, a ...any) (n int, err error) {
	return l.p(logLevel, fmt.Sprint(a...))
}

func (l *Logger) Progressf(format string, a ...any) (n int, err error) {
	if !l.isTerminal {
		return 0, nil
	}

	return l.p(LogLevelProgress, fmt.Sprintf(format, a...))
}

func (l *Logger) Printf(logLevel LogLevel, format string, a ...any) (n int, err error) {
	return l.p(logLevel, fmt.Sprintf(format, a...))
}

func (l *Logger) Println(logLevel LogLevel, a ...any) (n int, err error) {
	s := fmt.Sprintln(a...)
	return l.p(logLevel, s[:len(s)-1])
}

func (l *Logger) p(logLevel LogLevel, s string) (n int, err error) {
	if logLevel < l.Level {
		return 0, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := &msg{
		TimeStamp: l.timestamp(),
		Text:      s,
	}

	if l.StripMessages {
		msg.Text = strings.TrimSpace(msg.Text)
	}

	if l.isTerminal {
		if msg.TimeStamp != "" {
			msg.TimeStamp = l.TimeStampStyle.Render(msg.TimeStamp)
		}
		style, exists := l.Styles[logLevel]
		if exists && style != nil {
			msg.Text = l.Styles[logLevel].Render(msg.Text)
		}
	} else {
		msg.Prefix = l.prefix(logLevel)
	}

	str := msg.String()
	w := lipgloss.Width(str)

	if l.isTerminal && w < l.lastProgressLineWidth {
		str += strings.Repeat(" ", minNotLessThanZero(l.lastProgressLineWidth-w, l.GetWidth()-w))
		l.lastProgressLineWidth = 0
	}

	if logLevel == LogLevelProgress {
		l.lastProgressLineWidth = w
	}

	if logLevel == LogLevelProgress {
		str += "\r"
	} else {
		str += "\n"
	}

	return l.Writer.Write([]byte(str))
}
