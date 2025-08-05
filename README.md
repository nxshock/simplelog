# simplelog

Simple logging library for Go.

## Usage

```go
// Create logger
log := NewLogger(os.Stderr)     // To terminal - with colorful messages
log := NewLogger(any io.Writer) // To file - without colors

// Print messages
log.Info("info message")
log.Warnf("warning message: %s is suspicious", "something")
log.Errorf("error message: %s", "something goes wrong")
log.Fatal("unacceptable")

// Progress message
for i:=0; i < 100; i++ {
    log.Progressf("Processed %d records...", i)
}
log.Infof("Processed %d messages.")
```