package logger

import (
	"log"
)

var level int

// SetLevel is log level
func SetLevel(lv string) {
	switch lv {
	case "DEBUG":
		level = 0
	case "INFO":
		level = 1
	case "ERROR":
		level = 2
	default:
		level = 2
	}
}

// DEBUG is debug level log
func DEBUG(where string, what string, message string) {
	if level < 1 {
		write("DEBUG", where, what, message)
	}
}

// ERROR is error level log
func ERROR(where string, what string, message string) {
	write("ERROR", where, what, message)
}

// INFO is info level log
func INFO(where string, what string, message string) {
	if level < 2 {
		write("INFO", where, what, message)
	}
}

// FATAL is fatal error log
func FATAL(where string, what string, message string) {
	log.Fatalf("[Fatal ERROR] [%v -> %v] %v", where, what, message)
}

func write(level string, where string, what string, message string) {
	log.Printf("[%v] [%v -> %v] %v", level, where, what, message)
}
