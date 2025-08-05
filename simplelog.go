package simplelog

import (
	"cmp"
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
	LogLevelProgress LogLevel = iota
	LogLevelTrace
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelFatal
)

var (
	defaultFileTimestampFormat     = "2006-01-02 15:04:05"
	defaultTerminalTimestampFormat = "15:04:05"

	defaultTimestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultProgressStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultTraceStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultDebugStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))
	defaultInfoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc"))
	defaultWarningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffff80"))
	defaultErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	defaultFatalStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

type Logger struct {
	Writer io.Writer

	// timestamp format
	TimeFormat string

	// timestamp style
	TimeStampStyle lipgloss.Style

	// log level styles
	Styles map[LogLevel]lipgloss.Style

	// strip message from spaces before output
	StripMessages bool

	// is output to terminal
	isTerminal bool

	// last written progress message length
	lastProgressLineWidth int

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
		Styles:         make(map[LogLevel]lipgloss.Style),
		mu:             new(sync.Mutex)}

	logger.Styles[LogLevelProgress] = defaultProgressStyle
	logger.Styles[LogLevelTrace] = defaultTraceStyle
	logger.Styles[LogLevelDebug] = defaultDebugStyle
	logger.Styles[LogLevelInfo] = defaultInfoStyle
	logger.Styles[LogLevelWarning] = defaultWarningStyle
	logger.Styles[LogLevelError] = defaultErrorStyle
	logger.Styles[LogLevelFatal] = defaultFatalStyle

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

func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func levelSymbol(logLevel LogLevel) string {
	switch logLevel {
	case LogLevelTrace:
		return "TRC"
	case LogLevelDebug:
		return "DBG"
	case LogLevelInfo:
		return "INF"
	case LogLevelWarning:
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

func (l *Logger) Info(a ...any) (n int, err error) {
	return l.Print(LogLevelInfo, a...)
}

func (l *Logger) Infof(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelInfo, format, a...)
}

func (l *Logger) Warnf(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelWarning, format, a...)
}

func (l *Logger) Errorf(format string, a ...any) (n int, err error) {
	return l.Printf(LogLevelError, format, a...)
}

func (l *Logger) Infoln(a ...any) (n int, err error) {
	return l.Println(LogLevelInfo, a...)
}

func (l Logger) Fatal(a ...any) {
	l.mu.Lock()
	l.Print(LogLevelFatal, a...)
	l.mu.Unlock()

	os.Exit(1)
}

func (l *Logger) Fatalf(format string, a ...any) {
	l.Printf(LogLevelFatal, format, a...)

	os.Exit(1)
}

func (l *Logger) Fatalln(a ...any) {
	l.Println(LogLevelError, a...)

	os.Exit(1)
}

// clearAfterProgress clears line if next message length is shorter that previous message length.
func (l *Logger) clearAfterProgress(nextLogLevel LogLevel, nextWidth int) {
	if !l.isTerminal || l.lastProgressLineWidth == 0 || (nextLogLevel == LogLevelProgress && nextWidth >= l.lastProgressLineWidth) {
		return
	}

	fmt.Fprint(l.Writer, strings.Repeat(" ", min(l.GetWidth(), l.lastProgressLineWidth))+"\r")
	l.lastProgressLineWidth = 0
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
	return l.p(logLevel, fmt.Sprint(a...))
}

func (l *Logger) p(logLevel LogLevel, s string) (n int, err error) {
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
		msg.TimeStamp = l.TimeStampStyle.Render(msg.TimeStamp)
		msg.Text = l.Styles[logLevel].Render(msg.Text)
	} else {
		msg.Prefix = l.prefix(logLevel)
	}

	if logLevel == LogLevelProgress {
		msg.Text += "\r"
	} else {
		msg.Text += "\n"
	}

	str := msg.String()
	w := lipgloss.Width(str)
	l.clearAfterProgress(logLevel, w)

	if logLevel == LogLevelProgress {
		l.lastProgressLineWidth = w
	}

	return l.Writer.Write([]byte(str))
}
