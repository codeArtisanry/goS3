package s3

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	conf "github.com/codesnail21/gos3/config"
	"go.uber.org/zap"
)

type S3 struct {
	S3Client *s3.Client
	Cfg      conf.MinioConfig
}

// configS3 creates the S3 client
func New(minioConf conf.MinioConfig) (*S3, error) {
	endpointresolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{PartitionID: "aws",
			URL:               minioConf.Endpoint,
			SigningRegion:     minioConf.Region,
			HostnameImmutable: true}, nil
	})
	creds := credentials.NewStaticCredentialsProvider(minioConf.AccessKey, minioConf.SecretAccessKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), config.WithRegion(minioConf.Region), config.WithEndpointResolverWithOptions(endpointresolver))

	awsS3Client := s3.NewFromConfig(cfg)
	if err != nil {
		zap.L().Info("s3/main.go -> New  - ", zap.Error(err))
		return &S3{S3Client: awsS3Client, Cfg: minioConf}, err
	}
	return &S3{S3Client: awsS3Client, Cfg: minioConf}, nil
}

// Upload log file to Bucket
func (s3conn *S3) UploadFile(file io.Reader, filepath string) error {
	uploader := manager.NewUploader(s3conn.S3Client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s3conn.Cfg.Bucket),
		Key:    aws.String(filepath),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}

//Download file from S3 Bucket
func (s3conn *S3) DownloadS3File(filepath string) ([]byte, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})
	downloader := manager.NewDownloader(s3conn.S3Client)
	numBytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(s3conn.Cfg.Bucket),
		Key:    aws.String(filepath),
	})
	if err != nil {
		return nil, err
	}
	if numBytes < 1 {
		return nil, errors.New("zero bytes written to memory")
	}
	return buffer.Bytes(), nil
}
