package main

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func main() {
	client, err := oss.New("https://oss-cn-beijing.aliyuncs.com", "LTAI5t6kxoYpTPWZXw6ES6Gu", "hBUvanHY0OkbQG5IseB2KLzecypsjr")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	//// 列举所有的Buckets
	//listBuckets, err := client.ListBuckets()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	os.Exit(-1)
	//}
	//
	//for _, bucket := range listBuckets.Buckets {
	//	fmt.Println("Bucket:", bucket.Name)
	//}

	// 获取Bucket实例
	bucketName := "sample-douyin-video" // 请替换为你的Bucket名称
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error getting bucket:", err)
		os.Exit(-1)
	}

	// 列举Bucket中的所有文件
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			fmt.Println("Error listing objects:", err)
			os.Exit(-1)
		}

		// 打印文件名称
		for _, object := range lsRes.Objects {
			fmt.Println("Object:", object.Key)
		}

		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}
