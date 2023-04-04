//go:generate mockgen -source=logger.go -destination=../../mock/logger/logger.go -package=logger
/*
Package logger defines the interface for logging across all other packages.
*/
package logger

type Logger interface {
	Debugf(message string, args ...any)
	Infof(message string, args ...any)
	Warningf(message string, args ...any)
	Errorf(message string, args ...any)
	Fatalf(message string, args ...any)
}

// DummyLogger provides a default no-op logger
type DummyLogger struct{}

func (d DummyLogger) Debugf(message string, args ...any)   {}
func (d DummyLogger) Infof(message string, args ...any)    {}
func (d DummyLogger) Warningf(message string, args ...any) {}
func (d DummyLogger) Errorf(message string, args ...any)   {}
func (d DummyLogger) Fatalf(message string, args ...any)   {}
