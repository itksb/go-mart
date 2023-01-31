package logger

type Interface interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
