package main

import (
	"./storage"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

func GenSha256(filePath string) (s string, err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		err =  errors.New("Open metadata failed: " + err.Error())
		return
	}
	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h,f); err != nil {
		err = errors.New("Copy metadata to SHA256 handle failed: " + err.Error())
		return
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func UploadToBucket(checkPointFilePath string,
	fileName []string, filePath []string) (err error) {
	var oss storage.OSS
	err = oss.InitializeOss(Config.EndPoint,
		Config.AccessKeyID,
		Config.AccessKeySecret,
		Config.BucketName,
		checkPointFilePath)
	if err != nil {
		Log.Error("Oss initialize failed: ", err.Error())
		return
	}
	Log.Info("Oss initialize succeed")
	defer oss.Closer()

	if len(fileName) != len(filePath) {
		Log.Error("File name and path invalid. ")
		return errors.New("File name and path invalid. ")
	}

	for i := 0; i < len(fileName); i++ {
		Log.Debug("Uploading ", fileName[i])
		// Upload file to bucket root in default.
		err = oss.Upload(filePath[i], fileName[i])
		if err != nil { return }
	}

	return nil
}
