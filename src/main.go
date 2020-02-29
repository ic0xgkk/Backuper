package main

import (
	"flag"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "./config.json", "Config file path")
	flag.Parse()

	config, err := InitializeConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	var hostname string
	hostname, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	config.WorkDir += "/"
	globalLogFilePath := config.WorkDir + "Backuper-GlobalLog-" + hostname + ".log"
	var globalLogFile *os.File
	globalLogFile, err = os.OpenFile(globalLogFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}

	Log.Out = globalLogFile
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.TextFormatter{
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

	Log.Warn("Backuper starting...")
	Log.Info("Log is all set now and in Info level")
	Log.Info("Hostname: ", hostname)

	taskCron := gocron.NewScheduler()

	var task Task
	task.InitializeTask(config, hostname, taskCron)
	Log.Info("Task initialized")

	if config.StartTimeHour > 24 || config.StartTimeMinute > 60 {
		panic("What do you want to do? ")
	}
	taskCron.Every(uint64(config.PeriodDay)).Days().At(fmt.Sprint(config.StartTimeHour, ":", config.StartTimeMinute)).
		DoSafely(task.Start)
	Log.Info("Cron initialized")

	if config.ImmediateExec == true {
		Log.Info("Immediate execute enabled, processing")
		taskCron.RunAll()
	}

	<- taskCron.Start()

	_ = globalLogFile.Close()
}