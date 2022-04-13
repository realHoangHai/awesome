package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	root Logger
)

func init() {
	if root == nil {
		if err := Init(FromEnv()); err != nil {
			panic(err)
		}
	}
}

// Root return default logger instance.
// Root will try to init a root logger base on environment configuration.
// It will panic if failed to init.
func Root() Logger {
	if root == nil {
		if err := Init(FromEnv()); err != nil {
			panic(err)
		}
	}
	return root
}

// Init init the root logger with options.
func Init(opts ...Option) error {
	l := &Logh{}
	if err := l.Init(opts...); err != nil {
		return err
	}
	root = l
	return nil
}

// Infof print info with format.
func Infof(format string, v ...interface{}) {
	Root().Infof(format, v...)
}

// Debugf print debug with format.
func Debugf(format string, v ...interface{}) {
	Root().Debugf(format, v...)
}

// Warnf print warning with format.
func Warnf(format string, v ...interface{}) {
	Root().Warnf(format, v...)
}

// Errorf print error with format.
func Errorf(format string, v ...interface{}) {
	Root().Errorf(format, v...)
}

// Panicf panic with format.
func Panicf(format string, v ...interface{}) {
	Root().Panicf(format, v...)
}

// Info print info.
func Info(v ...interface{}) {
	Root().Info(v...)
}

// Debug print debug.
func Debug(v ...interface{}) {
	Root().Debug(v...)
}

// Warn print warning.
func Warn(v ...interface{}) {
	Root().Warn(v...)
}

// Error print error.
func Error(v ...interface{}) {
	Root().Error(v...)
}

// Panic panic.
func Panic(v ...interface{}) {
	Root().Panic(v...)
}

// Fields return a new logger entry with fields.
func Fields(kv ...interface{}) Logger {
	return Root().Fields(kv...)
}

// Context return a logger from the given context.
func Context(ctx context.Context) Logger {
	return Root().Context(ctx)
}

type (
	// Logh implement Logger interface using  logrus.
	Logh struct {
		logger *logrus.Entry
	}
)

// NewLogrus return new logger with context.
func NewLogger(opts ...Option) (*Logh, error) {
	l := &Logh{}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}
	return l, nil
}

// Init init the logger.
func (l *Logh) Init(opts ...Option) error {
	var f logrus.Formatter
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	format, err := options.GetFormat()
	if err != nil {
		return err
	}
	switch format {
	case FormatJSON:
		f = &logrus.JSONFormatter{
			TimestampFormat: options.TimeFormat,
		}
	case FormatText:
		f = &logrus.TextFormatter{
			TimestampFormat: time.RFC1123,
			FullTimestamp:   true,
		}
	}
	level, err := options.GetLevel()
	if err != nil {
		return err
	}
	out, err := options.GetWriter()
	if err != nil {
		return err
	}
	logger := logrus.New()
	logger.SetFormatter(f)
	logger.SetLevel(logrus.Level(level))
	logger.SetOutput(out)
	fields := map[string]interface{}{}
	for k, v := range options.Fields {
		fields[k] = v
	}
	l.logger = logrus.NewEntry(logger).WithFields(fields)
	return nil
}

// Info print info
func (l *Logh) Info(args ...interface{}) {
	l.logger.Infoln(args...)
}

// Debug print debug
func (l *Logh) Debug(v ...interface{}) {
	l.logger.Debugln(v...)
}

// Warn print warning
func (l *Logh) Warn(v ...interface{}) {
	l.logger.Warnln(v...)
}

// Error print error
func (l *Logh) Error(v ...interface{}) {
	l.logger.Errorln(v...)
}

// Panic panic
func (l *Logh) Panic(v ...interface{}) {
	l.logger.Panicln(v...)
}

// Infof print info with format.
func (l *Logh) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

// Debugf print debug with format.
func (l *Logh) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

// Warnf print warning with format.
func (l *Logh) Warnf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

// Errorf print error with format.
func (l *Logh) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(format, v...)
}

// Panicf panic with format.
func (l *Logh) Panicf(format string, v ...interface{}) {
	l.logger.Panicf(format, v...)
}

// Fields return a new logger with fields.
func (l *Logh) Fields(kv ...interface{}) Logger {
	return &Logh{
		logger: l.logger.WithFields(logrus.Fields(fields(kv...))),
	}
}

// Context return new logger from context.
func (l *Logh) Context(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		kv := make([]interface{}, 0)
		for k, v := range l.logger.Data {
			kv = append(kv, k, v)
		}
		return logger.Fields(kv...)
	}
	return l
}
