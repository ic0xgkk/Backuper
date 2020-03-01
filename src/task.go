package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"io"
	"os"
	"runtime"
	"time"
)

type Task struct {
	config             ConfigFile
	logFilePath        string
	finalFilePath      string
	checkPointFilePath string
	logFileName        string
	finalFileName      string
	globalLogFileName  string
	hostname           string
	scheduler          *gocron.Scheduler
}

func (t *Task) InitializeTask(c ConfigFile, hostname string, scheduler *gocron.Scheduler) () {
	t.config = c
	t.hostname = hostname
	t.scheduler = scheduler
}

func (t *Task) refreshTime() () {
	ts := time.Now()
	unixTs := fmt.Sprint(ts.Unix())
	Log.Info("Current time: ", ts)

	date := ts.Format("2006-01-02")
	t.logFileName = "Backuper-" + date + "-" + t.hostname + "-" + unixTs + ".log.gpg"
	t.finalFileName = "Backuper-" + date + "-" + t.hostname + "-" + unixTs + ".tar.gz.gpg"
	t.logFilePath = t.config.WorkDir + t.logFileName
	t.finalFilePath = t.config.WorkDir + t.finalFileName
	t.checkPointFilePath = t.config.WorkDir + "checkpoint.cp"

	Log.Info("Current task log file: ", t.logFileName)
	Log.Info("Current task backup file: ", t.finalFileName)
}

func (t *Task) start() (err error) {
	var logger Logger
	err = logger.Initialize(t.logFilePath, t.config.PubKeyPath)
	if err != nil { return }
	Log.Info("Backuper logger initialized")

	gnu := GnuPG{}
	err = gnu.InitializeGnuPG(t.config.PubKeyPath, t.finalFilePath)
	if err != nil { return }
	Log.Info("Backuper GnuPG initialized")

	compress := Compress{}
	err = compress.InitializeCompress(&logger, &gnu)
	if err != nil { return }
	Log.Info("Backuper compresser initialized")

	for i := 0; i < len(t.config.BackupPath); i++ {
		Log.Info("Processing ", t.config.BackupPath[i])
		err = compress.Compress(t.config.BackupPath[i])
		if err != nil { return }
	}
	Log.Info("Compress succeed")

	compress.Close()
	gnu.Close()

	logger.Info(fmt.Sprint("Backup created with ", logger.errCounter, " error"))
	var file *os.File
	file, err = os.Open(t.finalFilePath)
	if err != nil { return }
	h := sha256.New()
	if _, err = io.Copy(h,file); err != nil {
		return
	}
	_ = file.Close()
	logger.Info("Backup file sha256: ", fmt.Sprintf("%x", h.Sum(nil)))
	Log.Info("Logger succeed")

	logger.Close()

	oss := OSS{}
	err = oss.InitializeOss(t.config.EndPoint,
		t.config.AccessKeyID,
		t.config.AccessKeySecret,
		t.config.BucketName,
		t.checkPointFilePath)
	if err != nil { return }
	Log.Info("Oss initialized")

	Log.Info("Uploading ", t.logFileName)
	err = oss.Upload(t.logFilePath, t.logFileName)
	if err != nil { return }
	Log.Info("Succeed")

	Log.Info("Uploading ", t.finalFileName)
	err = oss.Upload(t.finalFilePath, t.finalFileName)
	if err != nil { return }
	Log.Info("Succeed")

	oss.Closer()

	if t.config.AutoDelete == true {
		err = os.Remove(t.logFilePath)
		if err != nil {
			Log.Warn("Auto delete failed, ", err.Error())
		}
		err = os.Remove(t.finalFilePath)
		if err != nil {
			Log.Warn("Auto delete failed, ", err.Error())
		}
	}

	return nil
}

func (t *Task) Start() () {
	Log.Warn("Task starting...")
	_, jt := t.scheduler.NextRun()
	Log.Info("Job current/next running time: ", jt)

	t.refreshTime()
	Log.Info("Time refreshed")

	err:= t.start()
	if err != nil {
		Log.Error("Backup failed, ", err.Error())
	}
	Log.Info("Backup succeed")

	runtime.GC()
}