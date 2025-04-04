package pkg

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gormlog "gorm.io/gorm/logger" // Importing the logger package from GORM
)

// GormLogger is a custom logger for GORM that uses zerolog for structured logging
// It implements the gorm.Logger interface for GORM to use it as a logger
type GormLogger struct {
	logger zerolog.Logger
}

// NewGormLogger creates a new GormLogger, call after InitLogger to set the format
// and log level for zerolog
func NewGormLogger() *GormLogger {
	return &GormLogger{
		logger: log.With().Str("component", "gorm").Logger(),
	}
}

// LogMode sets the log level for the GormLogger
func (l *GormLogger) LogMode(logLevel gormlog.LogLevel) gormlog.Interface {
	l.logger = l.logger.Level(zerolog.Level(toZeroLogLevel(logLevel)))
	return l
}

// Info logs an info message
func (l *GormLogger) Info(ctx context.Context, msg string, data ...any) {
	l.logger.Info().Ctx(ctx).Msgf(msg, data...)
}

// Warn logs a warning message
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...any) {
	l.logger.Warn().Ctx(ctx).Msgf(msg, data...)
}

// Error logs an error message
func (l *GormLogger) Error(ctx context.Context, msg string, data ...any) {
	l.logger.Error().Ctx(ctx).Msgf(msg, data...)
}

// Trace logs a trace message
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		l.logger.Trace().Ctx(ctx).Err(err).Msg("Trace error")
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()
	l.logger.Trace().Ctx(ctx).
		Dur("elapsed", elapsed).
		Int64("rows", rows).
		Str("sql", sql).
		Msg("Trace")
}

func InitLogger(format string) {
	zerolog.TimeFieldFormat = time.RFC3339
	switch format {
	case "pretty":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	default:
	}
}

func toZeroLogLevel(gormLevel gormlog.LogLevel) zerolog.Level {
	switch gormLevel {
	case gormlog.Silent:
		return zerolog.Disabled
	case gormlog.Error:
		return zerolog.ErrorLevel
	case gormlog.Warn:
		return zerolog.WarnLevel
	case gormlog.Info:
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}
