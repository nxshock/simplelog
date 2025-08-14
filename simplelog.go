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

	// Minimum log level of messages
	Level LogLevel

	// Marker of trimmed messages
	TrimMarker string

	// MinProgressUpdatePeriod is used to limit progress update speed
	MinProgressUpdatePeriod time.Duration

	// is output to terminal
	isTerminal bool

	// last written progress message length
	lastProgressLineWidth int

	// Timestamp of last written progress message
	lastProgressUpdateTime time.Time

	// mutex to prevent race conditions
	mu *sync.Mutex
}

// NewLogger returns new logger which writes messages to `w`.
func NewLogger(w io.Writer) *Logger {
	logger := &Logger{
		Writer:         w,
		TimeStampStyle: defaultTimestampStyle,
		Styles:         make(map[LogLevel]*lipgloss.Style),
		Level:          defaulLogLevel,
		TrimMarker:     defaultTrimMarker,
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

// getWidth returns current terminal width.
func (l *Logger) getWidth() int {
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

func (l *Logger) timestamp(t time.Time) string {
	if l.TimeFormat == "" {
		return ""
	}

	if !l.isTerminal {
		return t.Format(l.TimeFormat)
	}

	return l.TimeStampStyle.Render(t.Format(l.TimeFormat))
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
	timeStamp := time.Now()

	if logLevel == LogLevelProgress {
		if l.MinProgressUpdatePeriod > 0 && timeStamp.Sub(l.lastProgressUpdateTime) < l.MinProgressUpdatePeriod {
			return
		}

		l.lastProgressUpdateTime = timeStamp
	}

	if logLevel < l.Level {
		return 0, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := &msg{
		TimeStamp: l.timestamp(timeStamp),
		Text:      s,
	}

	if l.StripMessages {
		msg.Text = strings.TrimSpace(msg.Text)
	}

	if l.isTerminal {
		msg.fit(l.getWidth(), l.TrimMarker)

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
		str += strings.Repeat(" ", max(min(l.lastProgressLineWidth-w, l.getWidth()-w), 0))
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
