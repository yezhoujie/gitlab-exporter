package main

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func getOssBucket(config Config) *oss.Bucket {
	ossClient, err := oss.New(config.Oss.Endpoint, config.Oss.AccessKeyId, config.Oss.AccessSecret)
	if err != nil {
		log.Fatal("Failed to create oss client: ", err)
	}
	bucket, err := ossClient.Bucket(config.Oss.BucketName)
	if err != nil {
		log.Fatal("Failed to get bucket: ", err)
	}
	return bucket
}

// func getOssSign(bucket *oss.Bucket, objectName string) string {
// 	// 获取签名URL
// 	signUrl, err := bucket.SignURL(objectName, oss.HTTPPut, 3600)
// 	if err != nil {
// 		log.Fatal("Failed to get sign url: ", err)
// 	}
// 	// log.Println("signURL: ", signUrl)
// 	return signUrl
// }

func upload(bucket *oss.Bucket, objectName string) {
	// 上传文件
	err := bucket.PutObjectFromFile(objectName, objectName)
	if err != nil {
		log.Fatal("Failed to upload file: ", err)
	}
	log.Printf("file: %v upload to oss success, at bucket %v, path: %v\n", objectName, bucket.BucketName, objectName)
}
