package s3Client

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

type BucketConfig struct {
	AWS_ACCESS_KEY string
	AWS_SECRET_KEY string
	AWS_BUCKET     string
	AWS_REGION     string
}

type S3 struct {
	Session      *session.Session
	BucketConfig BucketConfig
	Uploader     *s3manager.Uploader
	Downloader   *s3manager.Downloader
}

func S3Client(bucketConfig BucketConfig) S3 {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	config := aws.Config{
		Region:      &bucketConfig.AWS_REGION,
		Credentials: credentials.NewStaticCredentials(bucketConfig.AWS_ACCESS_KEY, bucketConfig.AWS_SECRET_KEY, ""),
	}

	sess := session.Must(session.NewSession(&config))

	s3Client := S3{
		Session:      sess,
		BucketConfig: bucketConfig,
		Uploader:     s3manager.NewUploader(sess),
		Downloader:   s3manager.NewDownloader(sess),
	}
	return s3Client
}

func (s3Client S3) Download(filename string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := s3Client.Downloader.Download(buf, &s3.GetObjectInput{
		Bucket: &s3Client.BucketConfig.AWS_BUCKET,
		Key:    &filename,
	})
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to download file, %v", err)
	}
	return buf.Bytes(), nil
}

func (s3Client S3) Upload(filename string, contents interface{}) error {
	json, err := json.Marshal(contents)
	if err != nil {
		return fmt.Errorf("Failed to serialize struct, %v", err)
	}
	contents_bytes := bytes.NewReader(json)
	_, err = s3Client.Uploader.Upload(&s3manager.UploadInput{
		Bucket: &s3Client.BucketConfig.AWS_BUCKET,
		Key:    &filename,
		Body:   contents_bytes,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file, %v", err)
	}
	return nil
}
