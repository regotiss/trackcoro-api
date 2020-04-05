package objectstorage

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

var (
	sess     *session.Session
	svc      *s3.S3
	REGION   = "ap-south-1"
	S3Bucket = "trackcoro-images"
)

func InitializeS3Session() {
	logrus.Info("Initiating s3 session")
	cfg := aws.NewConfig().WithRegion(REGION)
	sess = session.Must(session.NewSession(cfg))
	svc = s3.New(sess)
	logrus.Info("S3 session initiated successfully")
}

func PutObject(key string, data []byte, bucket string) (*string, error) {
	if bucket == "" {
		bucket = S3Bucket
	}
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)
	size := fileBytes.Size()
	params := &s3.PutObjectInput{
		Bucket:        aws.String(S3Bucket),
		Key:           aws.String(key),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		logrus.Error("PutObject: ", err)
		return resp.ETag, err
	}
	return resp.ETag, err
}

func GetObject(key string, bucket string) ([]byte, error) {
	if bucket == "" {
		bucket = S3Bucket
	}
	object, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, object.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

