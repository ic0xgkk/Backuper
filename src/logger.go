package main

import (
	"fmt"
)

type Logger struct {
	fileName string
	gnupg    *GnuPG
	errCounter uint64
}

func (l *Logger) Initialize(logFilePath string, pubKeyPath string) (err error) {
	l.gnupg = &GnuPG{}
	err = l.gnupg.InitializeGnuPG(pubKeyPath, logFilePath)
	if err != nil { return }

	l.fileName = logFilePath
	l.errCounter = 0

	return nil
}

func (l *Logger) Info(msg ...interface{}) {
	message := fmt.Sprint("Info: ", msg, "\r\n")
	var err error
	_, err = l.gnupg.WriteCloser.Write([]byte(message))
	if err != nil {
		Log.Error("Writing a log failed: ", err.Error(), ". message: ", msg)
	}
}

func (l *Logger) Error(msg ...interface{}) {
	message := fmt.Sprint("Error: ", msg, "\r\n")
	var err error
	_, err = l.gnupg.WriteCloser.Write([]byte(message))
	if err != nil {
		Log.Error("Writing a log failed: ", err.Error(), ". message: ", msg)
	}
	l.errCounter++
}

func (l *Logger) Close() {
	l.gnupg.Close()
}
