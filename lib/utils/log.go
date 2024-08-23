package utils

import (
	"fmt"
	"time"
)

type LogLevel int

var STARTED = time.Now()

const (
	DEBUG LogLevel = 0
	INFO  LogLevel = 1
	WARN  LogLevel = 2
	ERROR LogLevel = 3
)

const LOG_LEVEL = INFO

func LogWithTaskId(id string, msg string, level LogLevel) {
	Since := time.Since(STARTED)
	if level >= LOG_LEVEL {
		fmt.Printf("%.1f [%s] %s\n", Since.Seconds(), id, msg)
	}
}
