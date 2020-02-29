package main

import (
	"fmt"
	"os"
)

type Logger struct {
	fileName string
	file *os.File
	errCounter uint64
}

func (l *Logger) Initialize(fileName string) (err error) {
	l.file, err = os.Create(fileName)
	if err != nil {
		return
	}
	l.fileName = fileName
	l.errCounter = 0
	return nil
}

func (l *Logger) Info(msg ...interface{}) {
	message := fmt.Sprint("Info: ", msg, "\r\n")
	var err error
	_, err = l.file.Write([]byte(message))
	if err != nil {
		Log.Error("Writing a log failed: ", err.Error(), ". message: ", msg)
	}
}

func (l *Logger) Error(msg ...interface{}) {
	message := fmt.Sprint("Error: ", msg, "\r\n")
	var err error
	_, err = l.file.Write([]byte(message))
	if err != nil {
		Log.Error("Writing a log failed: ", err.Error(), ". message: ", msg)
	}
	l.errCounter++
}

func (l *Logger) Close() {
	_ = l.file.Close()
}
