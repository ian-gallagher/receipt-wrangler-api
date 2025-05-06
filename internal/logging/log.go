package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
)

var logger *log.Logger
var stdLogger *log.Logger

func InitLog() error {
	logPath := "logs/app.log"
	logDir := "logs"
	if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(logDir, os.ModePerm)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	logFlags := log.Lshortfile | log.LstdFlags

	// Flags are for date time, file name, and line number
	logger = log.New(logFile, "App", logFlags)
	stdLogger = log.New(os.Stdout, "App", logFlags)

	return nil
}

func LogStd(level LogLevel, v ...any) {
	_, file, line, _ := runtime.Caller(1)
	lineInfo := fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
	levelString := fmt.Sprintf("%s: ", level)

	// Check if the first argument is an error
	var message string
	if len(v) > 0 {
		if err, ok := v[0].(error); ok {
			message = fmt.Sprintf("%s%s%s", lineInfo, levelString, err.Error())
		} else {
			message = fmt.Sprint(append([]any{lineInfo, levelString}, v...)...)
		}
	} else {
		message = fmt.Sprintf("%s%s", lineInfo, levelString)
	}

	if level == LOG_LEVEL_FATAL || level == LOG_LEVEL_ERROR {
		stackTrace := string(debug.Stack())
		message = fmt.Sprintf("%s\nStack Trace:\n%s", message, stackTrace)
	}

	logger.Println(message)
	stdLogger.Println(message)

	if level == LOG_LEVEL_FATAL {
		os.Exit(1)
	}
}
