package logger

import (
	"log"
	"os"
)

type Interface interface {
	Info(message string)
	Error(err error)
}

type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLog: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

func (l *Logger) Info(message string) {
	l.infoLog.Println(message)
}

func (l *Logger) Error(err error) {
	l.errorLog.Println(err)
}