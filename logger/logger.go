package logger

import (
	"fmt"
	"log/slog"
)

// Print log - for backward compatibility
func Print(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	msg := fmt.Sprint(v...)
	slog.Info(msg)
}

// Println log - for backward compatibility
func Println(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	msg := fmt.Sprint(v...)
	slog.Info(msg)
}

// Debug log
func Debug(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	
	// First argument is the message
	msg := fmt.Sprint(v[0])
	
	// If there are additional arguments, treat them as key-value pairs
	if len(v) > 1 {
		// Convert remaining arguments to key-value pairs
		args := convertArgsToKeyValues(v[1:])
		slog.Debug(msg, args...)
	} else {
		slog.Debug(msg)
	}
}

// Info log
func Info(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	
	// First argument is the message
	msg := fmt.Sprint(v[0])
	
	// If there are additional arguments, treat them as key-value pairs
	if len(v) > 1 {
		// Convert remaining arguments to key-value pairs
		args := convertArgsToKeyValues(v[1:])
		slog.Info(msg, args...)
	} else {
		slog.Info(msg)
	}
}

// Warn log
func Warn(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	
	// First argument is the message
	msg := fmt.Sprint(v[0])
	
	// If there are additional arguments, treat them as key-value pairs
	if len(v) > 1 {
		// Convert remaining arguments to key-value pairs
		args := convertArgsToKeyValues(v[1:])
		slog.Warn(msg, args...)
	} else {
		slog.Warn(msg)
	}
}

// Error log
func Error(v ...interface{}) {
	if len(v) == 0 {
		return
	}
	
	// First argument is the message
	msg := fmt.Sprint(v[0])
	
	// If there are additional arguments, treat them as key-value pairs
	if len(v) > 1 {
		// Convert remaining arguments to key-value pairs
		args := convertArgsToKeyValues(v[1:])
		slog.Error(msg, args...)
	} else {
		slog.Error(msg)
	}
}

// Enhanced API for structured logging

// DebugKV logs at Debug level with explicit key-value pairs
func DebugKV(msg string, keyValues ...interface{}) {
	slog.Debug(msg, keyValues...)
}

// InfoKV logs at Info level with explicit key-value pairs
func InfoKV(msg string, keyValues ...interface{}) {
	slog.Info(msg, keyValues...)
}

// WarnKV logs at Warn level with explicit key-value pairs
func WarnKV(msg string, keyValues ...interface{}) {
	slog.Warn(msg, keyValues...)
}

// ErrorKV logs at Error level with explicit key-value pairs
func ErrorKV(msg string, keyValues ...interface{}) {
	slog.Error(msg, keyValues...)
}

// Helper functions

// convertArgsToKeyValues converts a slice of interface{} to key-value pairs
// This helps maintain backward compatibility while enabling structured logging
func convertArgsToKeyValues(args []interface{}) []any {
	// For a single argument, add it as "details"
	if len(args) == 1 {
		return []any{"details", args[0]}
	}
	
	// For multiple arguments, try to pair them as key-value
	result := make([]any, 0, len(args))
	
	// If we have an odd number of arguments, the last one will be added as "details"
	for i := 0; i < len(args); i++ {
		if i+1 < len(args) {
			// Try to convert the key to string
			if key, ok := args[i].(string); ok {
				result = append(result, key, args[i+1])
				i++ // Skip the next item as we've used it as a value
			} else {
				// If key is not a string, add as a numbered item
				result = append(result, fmt.Sprintf("item_%d", i), args[i])
			}
		} else {
			// Last item with no pair
			result = append(result, fmt.Sprintf("item_%d", i), args[i])
		}
	}
	
	return result
}
