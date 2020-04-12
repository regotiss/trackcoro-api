package objectstorage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	svc      *s3.S3
	REGION   = "ap-south-1"
	S3Bucket = "trackcoro-images"
)

func InitializeS3Session() {
	logrus.Info("Initiating s3 session")
	cfg := aws.NewConfig().WithRegion(REGION)
	sess, err := session.NewSession(cfg)
	if err != nil {
		logrus.Panic("Could not initialize s3 session")
	}
	svc = s3.New(sess)
	logrus.Info("S3 session initiated successfully")
}

func PutObject(key string, data []byte, contentType string) (*string, error) {
	fileBytes := bytes.NewReader(data)
	size := fileBytes.Size()
	params := &s3.PutObjectInput{
		Bucket:        aws.String(S3Bucket),
		Key:           aws.String(key),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		logrus.Error("Could not save image: ", err)
		return resp.ETag, err
	}
	return resp.ETag, err
}

func GetObject(key string) ([]byte, error) {
	object, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logrus.Error("Could not download image: ", err)
		return nil, err
	}
	defer object.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, object.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

