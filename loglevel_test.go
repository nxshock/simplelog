package simplelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelMarshalText(t *testing.T) {
	for _, test := range []struct {
		logLevel LogLevel
		expected []byte
	}{
		{logLevel: LogLevelTrace, expected: []byte("trace")},
		{logLevel: LogLevelDebug, expected: []byte("debug")},
		{logLevel: LogLevelInfo, expected: []byte("info")},
		{logLevel: LogLevelWarn, expected: []byte("warn")},
		{logLevel: LogLevelError, expected: []byte("error")},
		{logLevel: LogLevelFatal, expected: []byte("fatal")},
	} {
		got, err := test.logLevel.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, test.expected, got)
	}
}

func TestLogLevelUnmarshalText(t *testing.T) {
	for _, test := range []struct {
		value    []byte
		expected LogLevel
	}{
		{value: []byte("trace"), expected: LogLevelTrace},
		{value: []byte("debug"), expected: LogLevelDebug},
		{value: []byte("info"), expected: LogLevelInfo},
		{value: []byte("warn"), expected: LogLevelWarn},
		{value: []byte("error"), expected: LogLevelError},
		{value: []byte("fatal"), expected: LogLevelFatal},
	} {
		var got LogLevel
		err := got.UnmarshalText(test.value)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, got)
	}
}
