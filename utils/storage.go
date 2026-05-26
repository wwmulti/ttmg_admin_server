package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// StorageDriver 定义统一的存储行为
type StorageDriver interface {
	Upload(content []byte, dir string, fileName string) (absoluteUrl string, relativePath string, err error)
	Delete(relativePath string) error
}

const (
	StorageLocal int = 1 // 本地存储
	StorageOSS   int = 2 // 阿里云OSS
	StorageAWS   int = 3 // AWS S3
)

// StorageFactory 工厂类
type StorageFactory struct{}

func (f *StorageFactory) GetDriver(storageType int) (StorageDriver, error) {
	switch storageType {
	case StorageLocal:
		return &LocalDriver{
			BaseDir: "public/static/uploads",
			Domain:  "http://localhost:8082",
		}, nil
	case StorageOSS:
		return &OssDriver{
			Endpoint:        "oss-cn-hangzhou.aliyuncs.com",
			AccessKeyId:     "你的ID",
			AccessKeySecret: "你的Secret",
			BucketName:      "你的Bucket",
			Domain:          "http://oss-cdn.example.com",
		}, nil
	/* case StorageAWS:
	return &AwsDriver{
		Region: "us-east-1",
		Bucket: "你的Bucket",
		Key:    "你的Key",
		Secret: "你的Secret",
		Domain: "http://s3-cdn.example.com",
	}, nil */
	default:
		return nil, fmt.Errorf("不支持的存储类型: %d", storageType)
	}
}

// --- LocalDriver 实现 ---
type LocalDriver struct {
	BaseDir string
	Domain  string
}

func (l *LocalDriver) Upload(content []byte, dir string, fileName string) (string, string, error) {
	relDir := filepath.Join(l.BaseDir, dir)
	if _, err := os.Stat(relDir); os.IsNotExist(err) {
		os.MkdirAll(relDir, os.ModePerm)
	}
	relPath := filepath.Join(relDir, fileName)
	if err := os.WriteFile(relPath, content, 0644); err != nil {
		return "", "", err
	}
	webPath := filepath.ToSlash(relPath)
	return strings.TrimSuffix(l.Domain, "/") + "/" + webPath, relPath, nil
}

func (l *LocalDriver) Delete(relPath string) error {
	return os.Remove(relPath)
}

// --- OssDriver 实现 ---
type OssDriver struct {
	Endpoint, AccessKeyId, AccessKeySecret, BucketName, Domain string
}

func (d *OssDriver) Upload(content []byte, dir string, fileName string) (string, string, error) {
	client, err := oss.New(d.Endpoint, d.AccessKeyId, d.AccessKeySecret)
	if err != nil {
		return "", "", err
	}
	bucket, err := client.Bucket(d.BucketName)
	if err != nil {
		return "", "", err
	}
	relPath := filepath.ToSlash(filepath.Join(dir, fileName))
	err = bucket.PutObject(relPath, bytes.NewReader(content))
	return strings.TrimSuffix(d.Domain, "/") + "/" + relPath, relPath, err
}

func (d *OssDriver) Delete(relPath string) error {
	client, _ := oss.New(d.Endpoint, d.AccessKeyId, d.AccessKeySecret)
	bucket, _ := client.Bucket(d.BucketName)
	return bucket.DeleteObject(relPath)
}

// --- AwsDriver 实现 (完全 V2 化) ---

/* type AwsDriver struct {
	Region, Bucket, Key, Secret, Domain string
}

func (d *AwsDriver) getS3Client() (*s3.Client, error) {
	// 加载配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(d.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(d.Key, d.Secret, "")),
	)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

func (d *AwsDriver) Upload(content []byte, dir string, fileName string) (string, string, error) {
	svc, err := d.getS3Client()
	if err != nil {
		return "", "", err
	}

	relPath := filepath.ToSlash(filepath.Join(dir, fileName))

	_, err = svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(relPath),
		Body:   bytes.NewReader(content),
	})

	if err != nil {
		return "", "", err
	}

	absUrl := strings.TrimSuffix(d.Domain, "/") + "/" + relPath
	return absUrl, relPath, nil
}

func (d *AwsDriver) Delete(relPath string) error {
	svc, err := d.getS3Client()
	if err != nil {
		return err
	}

	_, err = svc.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(relPath),
	})
	return err
} */
