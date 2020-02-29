package main

import (
	"encoding/json"
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
}

func InitializeConfig(configFilePath string) (fileConfig ConfigFile, err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(configFilePath); err != nil { return }

	var jsonByte []byte
	if jsonByte, err = ioutil.ReadAll(jsonFile); err != nil { return }

	if err = json.Unmarshal(jsonByte, &fileConfig); err != nil { return }

	_ = jsonFile.Close()

	return fileConfig, nil
}
