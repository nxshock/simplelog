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
estimated := uint(100)

progress := log.StartProgress(estimated) // start progress
for i:=0;i<100;i++ {
    progress.IncrementProgress(1) // incremeting progress will trigger updating of progress indicator
}
progress.Finish() // stop progress

log.Infof("Processed %d items.") // overwrite last progress message with your custom finish message
```