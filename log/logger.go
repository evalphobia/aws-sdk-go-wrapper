package log

// DefaultLogger is default Logger.
var DefaultLogger Logger

// Logger is logging interface.
type Logger interface {
	Infof(service, format string, v ...interface{})
	Errorf(service, format string, v ...interface{})
}

func init() {
	v := &DummyLogger{}
	DefaultLogger = v
}

// DummyLogger does not ouput anything
type DummyLogger struct{}

// Infof does nothing.
func (*DummyLogger) Infof(service, format string, v ...interface{}) {}

// Errorf does nothing.
func (*DummyLogger) Errorf(serivce, format string, v ...interface{}) {}
