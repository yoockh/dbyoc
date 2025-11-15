package logger

type Logger interface {
    Info(msg string, fields map[string]interface{})
    Warn(msg string, fields map[string]interface{})
    Error(msg string, fields map[string]interface{})
    Debug(msg string, fields map[string]interface{})
    WithFields(fields map[string]interface{}) Logger
}