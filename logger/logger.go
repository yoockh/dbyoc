package logger

import (
	"log"
	"os"
)

// Logger is an interface that defines the logging methods.
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

// SimpleLogger is a basic implementation of the Logger interface.
type SimpleLogger struct {
	logger *log.Logger
}

// NewSimpleLogger creates a new SimpleLogger instance.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info logs an informational message.
func (l *SimpleLogger) Info(msg string) {
	l.logger.Println("INFO: " + msg)
}

// Error logs an error message.
func (l *SimpleLogger) Error(msg string) {
	l.logger.Println("ERROR: " + msg)
}

// Debug logs a debug message.
func (l *SimpleLogger) Debug(msg string) {
	l.logger.Println("DEBUG: " + msg)
}

// LogQuery logs a database query.
func (l *SimpleLogger) LogQuery(query string) {
	l.Info("Executing query: " + query)
}

// LogError logs an error that occurred during a database operation.
func (l *SimpleLogger) LogError(err error) {
	l.Error("Database error: " + err.Error())
}