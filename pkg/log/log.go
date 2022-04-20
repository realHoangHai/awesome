package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
	"time"
)

type (
	// Logger defines standard operations of a logger.
	Logger interface {
		// Init init the logger.
		Init(...Option) error

		// Infof print info with format.
		Infof(format string, v ...interface{})

		// Debugf print debug with format.
		Debugf(format string, v ...interface{})

		// Warnf print warning with format.
		Warnf(format string, v ...interface{})

		// Errorf print error with format.
		Errorf(format string, v ...interface{})

		// Fatalf fatal with format.
		Fatalf(format string, v ...interface{})

		// Panicf panic with format.
		Panicf(format string, v ...interface{})

		// Info print info.
		Info(v ...interface{})

		// Debug print debug.
		Debug(v ...interface{})

		// Warn print warning.
		Warn(v ...interface{})

		// Error print error.
		Error(v ...interface{})

		// Fatal fatal.
		Fatal(v ...interface{})

		// Panic panic.
		Panic(v ...interface{})

		// Fields return new logger with the given fields.
		// The kv should be provided as key values pairs where key is a string.
		Fields(kv ...interface{}) Logger

		// Context provide a way to get a context logger,  i.e... with request-id.
		Context(ctx context.Context) Logger
	}

	// context key
	contextKey string

	// Options hold logger options
	Options struct {
		Level      Level
		Format     Format
		TimeFormat string
		Output     string
		Fields     map[string]string
		writer     io.Writer
	}
	// Option is an option for configure logger.
	Option = func(*Options)

	// Level is log level.
	Level int32

	// Format is log format
	Format string
)

const (
	loggerKey  contextKey = contextKey("logger_key")
	filePrefix            = "file://"
	// CorrelationID is field name of Correlation ID that is used to track related logs.
	CorrelationID string = "correlation_id"
)

// These are the different logging levels.
const (
	// LevelPanic level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic Level = iota
	// LevelFatal level. Logs and then calls os.Exit. It will exit even if the
	// logging level is set to Panic.
	LevelFatal
	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError
	// LevelWarn level. Non-critical entries that deserve eyes.
	LevelWarn
	// LevelInfo level. General operational entries about what's going on inside the
	// application.
	LevelInfo
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug
	// LevelTrace level. Designates finer-grained informational events than the Debug.
	LevelTrace
)

// Formats of log output.
const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// NewContext return a new logger context.
func NewContext(ctx context.Context, logger Logger) context.Context {
	if logger == nil {
		logger = Root()
	}
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext get logger form context.
func FromContext(ctx context.Context) Logger {
	if ctx == nil {
		return Root()
	}
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return Root()
}

// GetWriter return writer output. If the given output is not valid, os.Stdout is returned.
func (opts Options) GetWriter() (io.Writer, error) {
	switch {
	case opts.writer != nil:
		return opts.writer, nil
	case strings.HasPrefix(opts.Output, filePrefix):
		name := opts.Output[len(filePrefix):]
		f, err := os.Create(name)
		if err != nil {
			return nil, err
		}
		return f, nil
	case opts.Output == "":
		return os.Stderr, nil
	default:
		return nil, fmt.Errorf("log: output not supported: %s", opts.Output)
	}
}

// GetLevel return log level. If the given level is not valid, LevelDebug is returned.
func (opts Options) GetLevel() (Level, error) {
	if opts.Level < LevelPanic || opts.Level > LevelTrace {
		return LevelDebug, fmt.Errorf("log: level not supported: %d", opts.Level)
	}
	return opts.Level, nil
}

// GetFormat return format of output log. If given format is not valid, JSON format is returned.
func (opts Options) GetFormat() (Format, error) {
	if opts.Format != FormatText && opts.Format != FormatJSON && opts.Format != "" {
		return "", fmt.Errorf("log: format not supported: %s", opts.Format)
	}
	if opts.Format == "" {
		return FormatJSON, nil
	}
	return opts.Format, nil
}

func fields(kv ...interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	ood := len(kv) % 2
	for i := 0; i < len(kv)-ood; i += 2 {
		fields[fmt.Sprintf("%v", kv[i])] = kv[i+1]
	}
	if ood == 1 {
		fields["msg.1"] = fmt.Sprintf("%v", kv[len(kv)-1])
	}
	return fields
}

// FromOptions is an option to create new logger from an existing Options.
func FromOptions(opts *Options) Option {
	return func(v *Options) {
		v.Fields = opts.Fields
		v.Format = opts.Format
		v.Level = opts.Level
		v.TimeFormat = opts.TimeFormat
	}
}

// WithLevel provides an option to set log level.
func WithLevel(level Level) Option {
	return func(opts *Options) {
		opts.Level = level
	}
}

// WithFormat provides an option to set log format.
func WithFormat(f Format) Option {
	return func(opts *Options) {
		opts.Format = f
	}
}

// WithTimeFormat provides an option to set time format for logger.
func WithTimeFormat(f string) Option {
	return func(opts *Options) {
		opts.TimeFormat = f
	}
}

// WithFields provides an option to set logger fields.
func WithFields(kv ...interface{}) Option {
	return func(opts *Options) {
		if opts.Fields == nil {
			opts.Fields = make(map[string]string)
		}
		for k, v := range fields(kv...) {
			opts.Fields[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
		}
	}
}

// WithWriter provides an option to set a output writer.
func WithWriter(w io.Writer) Option {
	return func(opts *Options) {
		opts.writer = w
	}
}

// Logh implements Logger interface using logrus.
type Logh struct {
	logger *logrus.Entry
}

func NewLogh(opts ...Option) (*Logh, error) {
	l := &Logh{}
	if err := l.Init(opts...); err != nil {
		return nil, err
	}
	return l, nil
}

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

func (l *Logh) Info(v ...interface{}) {
	l.logger.Infoln(v...)
}

func (l *Logh) Debug(v ...interface{}) {
	l.logger.Debugln(v...)
}

func (l *Logh) Warn(v ...interface{}) {
	l.logger.Warnln(v...)
}

// Error print error
func (l *Logh) Error(v ...interface{}) {
	l.logger.Errorln(v...)
}

func (l *Logh) Fatal(v ...interface{}) {
	l.logger.Fatalln(v...)
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

// Fatalf fatal with format.
func (l *Logh) Fatalf(format string, v ...interface{}) {
	l.logger.Fatalf(format, v...)
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
