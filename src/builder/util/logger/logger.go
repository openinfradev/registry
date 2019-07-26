package logger

import (
	"log"
)

var level int

// SetLevel is log level
func SetLevel(lv int) {
	level = lv
}

// DEBUG is debug level log
func DEBUG(where string, message string) {
	if level < 1 {
		write("DEBUG", where, message)
	}
}

// ERROR is error level log
func ERROR(where string, message string) {
	write("ERROR", where, message)
}

// INFO is info level log
func INFO(where string, message string) {
	if level < 2 {
		write("INFO", where, message)
	}
}

// FATAL is fatal error log
func FATAL(where string, message string) {
	log.Fatalf("[Fatal ERROR] [%v] %v", where, message)
}

func write(level string, where string, message string) {
	log.Printf("[%v] [%v] %v", level, where, message)
}
