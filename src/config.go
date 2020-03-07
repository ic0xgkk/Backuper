package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
	WorkDir         string   `json:"work_dir"`
	PubKeyPath      string   `json:"pub_key_path"`
	EndPoint        string   `json:"end_point"`
	AccessKeyID     string   `json:"access_key_id"`
	AccessKeySecret string   `json:"access_key_secret"`
	BucketName      string   `json:"bucket_name"`
	BackupPath      []string `json:"backup_path"`
	PeriodDay       uint32   `json:"period_day"`
	StartTimeHour   uint8    `json:"start_time_hour"`
	StartTimeMinute uint8    `json:"start_time_minute"`
	AutoDelete      bool     `json:"auto_delete"`
	ImmediateExec   bool     `json:"immediate_exec"`
	LogLevel        string   `json:"log_level"`
}

var Config *ConfigFile

func InitializeConfig(configFilePath string) (err error) {
	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		return errors.New("Open global log file failed: " + err.Error())
	}
	defer jsonFile.Close()

	var jsonByte []byte
	jsonByte, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.New("Read json config failed: " + err.Error())
	}

	err = json.Unmarshal(jsonByte, &Config)
	if err != nil {
		return errors.New("Unmarshal json failed: " + err.Error())
	}

	return nil
}
