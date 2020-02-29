package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"runtime"
)

type OSS struct {
	client *oss.Client
	bucket *oss.Bucket
	checkPointFilePath string
}

func (o *OSS) InitializeOss(endPoint string,
	accessKeyID string,
	accessKeySecret string,
	bucketName string,
	checkPointFilePath string) (err error) {
	o.client, err = oss.New(endPoint, accessKeyID, accessKeySecret)
	if err != nil {
		return
	}

	o.bucket, err = o.client.Bucket(bucketName)
	if err != nil {
		return
	}

	o.checkPointFilePath = checkPointFilePath

	return nil
}

func (o *OSS) Upload(localFilePath string, objectName string) (err error) {
	err = o.bucket.UploadFile(objectName, localFilePath, 10*1024*1024, oss.Routines(5),
		oss.Checkpoint(true, o.checkPointFilePath))
	if err != nil {
		return
	}
	return nil
}

func (o *OSS) Closer()  {
	o.client = nil
	o.bucket = nil
	runtime.GC()
}