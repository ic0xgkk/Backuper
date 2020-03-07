package main

import (
	"./compress"
	"./encrypt"
	"./file"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"os"
	"runtime"
	"time"
)

type Task struct {
	logFilePath        string
	finalFilePath      string
	checkPointFilePath string

	logFileName        string
	finalFileName      string
	scheduler          *gocron.Scheduler

	partFile           file.File
	partFileEncrypt    encrypt.Encrypt

	metaFile           file.File
	metaFileEncrypt    encrypt.Encrypt
}

func (t *Task) Initialize(scheduler *gocron.Scheduler) () {
	t.scheduler = scheduler
}

func (t *Task) refreshTime() () {
	ts := time.Now()
	unixTs := fmt.Sprint(ts.Unix())
	Log.Info("Current time: ", ts)

	date := ts.Format("2006-01-02")
	t.logFileName = "Backuper-" + date + "-" + Env.Hostname + "-" + unixTs + ".log.gpg"
	t.finalFileName = "Backuper-" + date + "-" + Env.Hostname + "-" + unixTs + ".tar.gz.gpg"
	t.logFilePath = Config.WorkDir + "/" + t.logFileName
	t.finalFilePath = Config.WorkDir + "/" + t.finalFileName
	t.checkPointFilePath = Config.WorkDir + "/checkpoint.cp"

	Log.Info("Current task log file: ", t.logFileName)
	Log.Info("Current task backup file: ", t.finalFileName)
}

func (t *Task) TurnOnMultiLogger() (err error) {
	err = t.partFile.Initialize(t.logFilePath, "rw")
	if err != nil {
		Log.Error("Open part log file failed: ", err.Error())
		return 
	}

	err = t.partFileEncrypt.Initialize(t.partFile.GetWriter(), Config.PubKeyPath)
	if err != nil {
		Log.Error("Failed to initialize part encrypted log: ", err.Error())
		return
	}

	Log.SwitchToMultiWriter(t.partFileEncrypt.GetWriter())

	return nil
}

func (t *Task) TurnOffMultiLogger() {
	Log.SwitchToSingleWriter()
	t.partFileEncrypt.Close()
	t.partFile.Close()
}

func (t *Task) PrepareMetadata() (err error) {
	err = t.metaFile.Initialize(t.finalFilePath, "rw")
	if err != nil {
		Log.Error("Initialize metadata file failed: ", err.Error())
		return
	}

	err = t.metaFileEncrypt.Initialize(t.metaFile.GetWriter(), Config.PubKeyPath)
	if err != nil {
		Log.Error("Initialize metadata file encrypt failed: ", err.Error())
		return
	}

	return nil
}

func (t *Task) ClearMetadata() {
	t.metaFileEncrypt.Close()
	t.metaFile.Close()
}

func (t *Task) Start() () {
	Log.Warn("Task starting...")

	_, jt := t.scheduler.NextRun()
	Log.Debug("Job current/next running time: ", jt)

	t.refreshTime()
	Log.Info("Time refreshed")

	err := t.PrepareMetadata()
	if err != nil {
		Log.Error("Prepare metadata failed")
		return
	}
	Log.Info("Prepare metadata succeed")
	defer t.ClearMetadata()

	err = t.TurnOnMultiLogger()
	if err != nil {
		Log.Error("Turn on multi logger failed")
		return
	}
	Log.Info("Turn on multi logger")
	defer t.TurnOffMultiLogger()

	var comp compress.Compress
	err = comp.Initialize(t.metaFileEncrypt.GetWriter(), FileInfoHook)
	if err != nil {
		Log.Error("Initialize compress failed: ", err.Error())
		return
	}
	Log.Info("Initialize compress succeed")
	defer comp.Close()

	for i := 0; i < len(Config.BackupPath); i++ {
		Log.Info("Processing ", Config.BackupPath[i])
		err = comp.Compress(Config.BackupPath[i])
		if err != nil {
			Log.Error("Compress failed: ", err.Error())
			return
		}
	}
	comp.Close()  // Write compress footer manually.
	Log.Info("Compress handle closed")

	Log.Info("Compress succeed")

	t.ClearMetadata()  // Write encrypt footer and file footer manually.
	Log.Info("Metadata handle closed")

	var sha256code string
	sha256code, err = GenSha256(t.finalFilePath)
	if err != nil {
		Log.Error("Generate final metadata sha256 failed: ", err.Error())
		return
	}
	Log.Info("Metadata sha256 code: ", sha256code)

	t.TurnOffMultiLogger()  // Switch to single logger mode and turn off file descriptor manually.
	Log.Info("Turn off multi logger")

	var names, paths []string
	err = UploadToBucket(t.checkPointFilePath,
		append(names, t.finalFileName, t.logFileName),
		append(paths, t.finalFilePath, t.logFilePath))
	if err != nil {
		Log.Error("Upload failed")
	} else {
		Log.Info("Upload succeed")
	}

	if Config.AutoDelete == true {
		Log.Info("Auto delete enabled")
		err = os.Remove(t.logFilePath)
		if err != nil {
			Log.Warn("Auto delete failed, ", err.Error())
		}
		err = os.Remove(t.finalFilePath)
		if err != nil {
			Log.Warn("Auto delete failed, ", err.Error())
		}
	}

	Log.Warn("Task finished")

	runtime.GC()
}