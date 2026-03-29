package simplelog

import (
	"fmt"
	"strings"
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

var stringToLogLevel = map[string]LogLevel{
	"trace": LogLevelTrace,
	"debug": LogLevelDebug,
	"info":  LogLevelInfo,
	"warn":  LogLevelWarn,
	"error": LogLevelError,
	"fatal": LogLevelFatal}

var logLevelToString = map[LogLevel]string{
	LogLevelTrace: "trace",
	LogLevelDebug: "debug",
	LogLevelInfo:  "info",
	LogLevelWarn:  "warn",
	LogLevelError: "error",
	LogLevelFatal: "fatal"}

func (l *LogLevel) UnmarshalText(text []byte) error {
	v, ok := stringToLogLevel[strings.ToLower(string(text))]
	if !ok {
		return fmt.Errorf("unknown value: %s", string(text))
	}

	*l = v

	return nil
}

func (l *LogLevel) MarshalText() (text []byte, err error) {
	return []byte(logLevelToString[LogLevel(*l)]), nil
}
