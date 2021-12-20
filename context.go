package log

import (
	"github.com/wspowell/context"
)

type configContextKey struct{}
type levelContextKey struct{}
type logContextKey struct{}

// WithConfig adds a Configer to the Context which enables new Loggers to be created.
// This creates a new Logger from the Configer and adds it as a local Context value.
func WithContext(ctx context.Context, config Configuration) context.Context {
	ctx = context.WithValue(ctx, configContextKey{}, config)
	ctx = context.WithValue(ctx, levelContextKey{}, config.Level())
	withLogger(ctx, config.Logger())

	return ctx
}

func withLogger(ctx context.Context, logger Logger) {
	// Logger is localized to the context so it guarantees a data race cannot occur.
	context.WithLocalValue(ctx, logContextKey{}, logger)
}

func fromContext(ctx context.Context, logLevel Level) (Logger, bool) {
	if configLevel, ok := ctx.Value(levelContextKey{}).(Level); ok {
		if logLevel < configLevel {
			return nil, false
		}

		return getLogger(ctx)
	}

	return nil, false
}

func getLogger(ctx context.Context) (Logger, bool) {
	if log, ok := ctx.Value(logContextKey{}).(Logger); ok {
		return log, true
	}

	return newLogger(ctx)
}

func newLogger(ctx context.Context) (Logger, bool) {
	// Create new Logger from the Configer, if present.
	if config, ok := ctx.Value(configContextKey{}).(Configuration); ok {
		log := config.Logger()
		withLogger(ctx, log)

		return log, true
	}

	return nil, false
}

func Tag(ctx context.Context, name string, value any) {
	if log, ok := fromContext(ctx, LevelPanic); ok {
		log.Tag(name, value)
	}
}

func Tags(ctx context.Context) map[string]any {
	if log, ok := fromContext(ctx, LevelPanic); ok {
		return log.Tags()
	}

	return nil
}

func Printf(ctx context.Context, format string, v ...any) {
	// Log at INFO to match logrus.
	if log, ok := fromContext(ctx, LevelInfo); ok {
		log.Printf(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Trace(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelTrace); ok {
		log.Trace(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Debug(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelDebug); ok {
		log.Debug(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Info(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelInfo); ok {
		log.Info(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Warn(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelWarn); ok {
		log.Warn(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Error(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelError); ok {
		log.Error(format, v...)
	}
}

// nolint:goprintffuncname // reason: keep in line with logger function naming
func Fatal(ctx context.Context, format string, v ...any) {
	if log, ok := fromContext(ctx, LevelFatal); ok {
		log.Fatal(format, v...)
	}
}
