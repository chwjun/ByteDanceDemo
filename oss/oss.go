package oss

import (
	"fmt"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client
var Bucket *oss.Bucket

const (
	endpoint        = "https://oss-cn-beijing.aliyuncs.com"
	endpointname    = "oss-cn-beijing.aliyuncs.com"
	accessID        = "LTAI5t6kxoYpTPWZXw6ES6Gu"
	accessKey       = "hBUvanHY0OkbQG5IseB2KLzecypsjr"
	bucketName      = "sample-douyin-video"
	URLPre          = "https://" + bucketName + "." + endpointname
	CoverURL_SUFFIX = "?x-oss-process=video/snapshot,t_1000,m_fast"
)

func Init() {
	//创建实例
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Printf("client:%#v\n", client)
	// err = client.CreateBucket(bucketName)
	isExist, err := client.IsBucketExist(bucketName)
	if err != nil {
		log.Println("IsBucketExist Error:", err)
	}
	// 没有该存储空间
	if !isExist {
		err = client.CreateBucket(bucketName, oss.StorageClass(oss.StorageIA), oss.ACL(oss.ACLPublicReadWrite), oss.RedundancyType(oss.RedundancyZRS))
		if err != nil {
			log.Println("Create Bucket ERROR : ", err)

		}
	}
	Bucket, err = client.Bucket(bucketName)
	if err != nil {
		log.Println("GetBucket ERROR : ", err)
	}
}
