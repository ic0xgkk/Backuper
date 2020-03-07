package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type Logger struct {
	logrus.Logger
	globalWriter io.Writer
	partWriter   io.Writer
}

func (l *Logger) Close() {
	l.Close()
}

func (l *Logger) Initialize(level string, file *os.File)  {
	switch strings.ToLower(level) {
	default:
	case "debug":
		l.SetLevel(logrus.DebugLevel)
		break
	case "info":
		l.SetLevel(logrus.InfoLevel)
		break
	case "warn":
	case "warning":
		l.SetLevel(logrus.WarnLevel)
		break
	case "error":
		l.SetLevel(logrus.ErrorLevel)
		break
	}

	l.globalWriter = io.MultiWriter(file)

	l.Out = l.globalWriter

	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             false,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	})
	l.Info("Log level: ", strings.ToLower(level))
	l.Info("Logger initialize succeed")
}


func (l *Logger) SwitchToMultiWriter(writer *io.WriteCloser) {
	l.Out = io.MultiWriter(l.globalWriter, *writer)
}

func (l *Logger) SwitchToSingleWriter() {
	l.Out = l.globalWriter
}

