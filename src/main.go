package main

import (
	"./file"
	"./logger"
	"flag"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
)

type Environment struct {
	Hostname string
}

var Env Environment
var Log logger.Logger

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "./config.json", "Config file path")
	flag.Parse()

	err := InitializeConfig(configFilePath)
	if err != nil {
		panic(err.Error())
	}
	if Config.StartTimeHour > 24 || Config.StartTimeMinute > 60 {
		panic("StartTimeHour or StartTimeMinute error")
	}
	if Config.PeriodDay == 0 {
		panic("PeriodDay can not be zero")
	}

	Env.Hostname, err = os.Hostname()
	if err != nil {
		panic("Get hostname failed: " + err.Error())
	}

	var logFile file.File
	err = logFile.Initialize(Config.WorkDir + "/Backuper-GlobalLog-" + Env.Hostname + ".log", "wa")
	if err != nil {
		panic("Initialize global log failed: " + err.Error())
	}
	defer logFile.Close()
	Log.Initialize(Config.LogLevel, logFile.GetWriter())
	Log.Info("Hostname: ", Env.Hostname)

	taskCron := gocron.NewScheduler()

	var task Task
	task.Initialize(taskCron)
	taskCron.Every(uint64(Config.PeriodDay)).Days().At(fmt.Sprint(Config.StartTimeHour, ":", Config.StartTimeMinute)).
		DoSafely(task.Start)
	Log.Info("Cron initialized")

	if Config.ImmediateExec == true {
		Log.Info("Immediate execute enabled, processing now")
		taskCron.RunAll()
	}

	<- taskCron.Start()
}

