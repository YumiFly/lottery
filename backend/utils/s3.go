// utils/s3.go
package utils

import (
	"backend/config"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

var S3Client *s3.Client

// InitS3Client 初始化 S3 客户端
func InitS3Client() {
	// 检查必要的配置是否存在
	if config.AppConfig.AccessKey == "" || config.AppConfig.SecretKey == "" ||
		config.AppConfig.Region == "" || config.AppConfig.BucketName == "" {
		Logger.Warning("S3 configuration is incomplete, S3 storage will not be available")
		return
	}

	// 创建 AWS 凭证
	creds := credentials.NewStaticCredentialsProvider(
		config.AppConfig.AccessKey,
		config.AppConfig.SecretKey,
		"",
	)

	// 创建 S3 客户端
	options := s3.Options{
		Region:      config.AppConfig.Region,
		Credentials: creds,
	}

	// 如果指定了自定义端点，则使用自定义端点
	if config.AppConfig.Endpoint != "" {
		// 使用新的 BaseEndpoint 方式设置端点
		options.BaseEndpoint = aws.String(config.AppConfig.Endpoint)
		// 如果使用自定义端点，可能需要禁用虚拟主机匹配
		options.UsePathStyle = true
	}

	S3Client = s3.New(options)

	Logger.Info("S3 client initialized successfully")
}

// UploadFileToS3 将文件上传到 S3
func UploadFileToS3(file *multipart.FileHeader, directory string) (string, error) {
	if S3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 读取文件内容
	buffer := make([]byte, file.Size)
	_, err = src.Read(buffer)
	if err != nil {
		return "", err
	}

	// 生成唯一的文件名
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	key := filepath.Join(directory, filename)

	// 上传到 S3
	_, err = S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(config.AppConfig.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(getContentType(file.Filename)),
	})

	if err != nil {
		Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"file":  file.Filename,
		}).Error("Failed to upload file to S3")
		return "", err
	}

	// 构造文件 URL
	var fileURL string
	if config.AppConfig.Endpoint != "" {
		// 使用自定义端点
		fileURL = fmt.Sprintf("%s/%s/%s",
			config.AppConfig.Endpoint,
			config.AppConfig.BucketName,
			key)
	} else {
		// 使用默认的 AWS S3 URL 格式
		fileURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
			config.AppConfig.BucketName,
			config.AppConfig.Region,
			key)
	}

	Logger.WithFields(logrus.Fields{
		"file": file.Filename,
		"url":  fileURL,
	}).Info("File uploaded to S3 successfully")

	return fileURL, nil
}

// 根据文件扩展名获取内容类型
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// GetFileFromS3 从 S3 获取文件
func GetFileFromS3(key string) ([]byte, error) {
	if S3Client == nil {
		return nil, fmt.Errorf("S3 client not initialized")
	}

	// 从 S3 获取对象
	result, err := S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(config.AppConfig.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("Failed to get file from S3")
		return nil, err
	}
	defer result.Body.Close()

	// 读取文件内容
	return io.ReadAll(result.Body)
}

// DeleteFileFromS3 从 S3 删除文件
func DeleteFileFromS3(key string) error {
	if S3Client == nil {
		return fmt.Errorf("S3 client not initialized")
	}

	// 从 S3 删除对象
	_, err := S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(config.AppConfig.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("Failed to delete file from S3")
		return err
	}

	Logger.WithFields(logrus.Fields{
		"key": key,
	}).Info("File deleted from S3 successfully")

	return nil
}
