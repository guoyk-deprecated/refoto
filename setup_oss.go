package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"path"
	"strings"
	"time"
)

var (
	ossSuffixAvatar  = "refoto_avatar"
	ossSuffixPreview = "refoto_preview"
)

func setupOSS() (bucket *oss.Bucket, err error) {
	var client *oss.Client
	if client, err = oss.New(envOSSEndpoint, envOSSAccessKeyID, envOSSAccessKeySecret); err != nil {
		return
	}
	if bucket, err = client.Bucket(envOSSBucket); err != nil {
		return
	}
	return
}

func ossUploadFile(bucket *oss.Bucket, filename string, file io.Reader) (relPath string, err error) {
	now := time.Now()
	relPath = path.Join(
		"refoto",
		"avatar",
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%d", now.Unix())+"-"+path.Base(filename),
	)
	if err = bucket.PutObject(relPath, file); err != nil {
		return
	}
	return
}

func ossCombineURL(relPath string, suffix string) string {
	if suffix == "" {
		return strings.TrimSuffix(envOSSPublicEndpoint, "/") + "/" + strings.TrimPrefix(relPath, "/")
	}
	return strings.TrimSuffix(envOSSPublicEndpoint, "/") + "/" + strings.TrimPrefix(relPath, "/") + "!" + suffix
}
