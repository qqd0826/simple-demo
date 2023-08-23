package db

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"
	//"strconv"
	//"time"
)

var bucket *oss.Bucket

func InitOss() {
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	endpoint := "oss-cn-hangzhou.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := ""
	accessKeySecret := ""
	// yourBucketName填写存储空间名称。
	bucketName := ""
	// uploadFileName填写文件上传的位置及名字。
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
}
func UploadHandler(UserId int64, title string, file *multipart.FileHeader) (uploadPath string, err error) {
	f, err := file.Open()
	if err != nil {
		log.Fatal("接收文件失败", err)
		return "", err
	}
	uploadFileName := strconv.FormatInt(UserId, 10) + "-" + title + "-" + strconv.FormatInt(time.Now().Unix(), 10) + path.Ext(file.Filename)
	fmt.Println(uploadFileName)
	// 上传文件。
	err = bucket.PutObject(uploadFileName, f)
	if err != nil {
		handleError(err)
		return "", err
	}
	return uploadFileName, nil
}
func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}
