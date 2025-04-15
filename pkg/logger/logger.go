package logger

import (
	"time"

	"github.com/lamasutra/bg-music/pkg/events"
)

const (
	LevelFatal = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

type logBuffer struct {
	level    uint8
	messages []message
}

var defaultLogger *logBuffer = &logBuffer{
	level:    LevelTrace,
	messages: make([]message, 0),
}

func New(level uint8) *logBuffer {
	return &logBuffer{
		messages: make([]message, 0),
	}
}

func SetLevel(level uint8) {
	defaultLogger.level = level
}

func (lb *logBuffer) Log(level uint8, args ...any) {
	if level < lb.level {
		return
	}
	msg := message{
		level: level,
		time:  time.Now(),
		msg:   args,
	}
	lb.messages = append(lb.messages, msg)
	events.Trigger("log", msg)
}

func (lb *logBuffer) Fatal(args ...any) {
	lb.Log(LevelFatal, args...)
}

func (lb *logBuffer) Error(args ...any) {
	lb.Log(LevelError, args...)
}

func (lb *logBuffer) Warn(args ...any) {
	lb.Log(LevelWarn, args...)
}

func (lb *logBuffer) Info(args ...any) {
	lb.Log(LevelInfo, args...)
}

func (lb *logBuffer) Debug(args ...any) {
	lb.Log(LevelDebug, args...)
}

func (lb *logBuffer) Trace(args ...any) {
	lb.Log(LevelTrace, args...)
}

func Log(level uint8, args ...any) {
	defaultLogger.Log(level, args...)
}

func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

func Error(args ...any) {
	defaultLogger.Error(args...)
}

func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

func Info(args ...any) {
	defaultLogger.Info(args...)
}

func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

func Trace(args ...any) {
	defaultLogger.Trace(args...)
}
