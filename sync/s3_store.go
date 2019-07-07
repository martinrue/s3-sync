package sync

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Store is an implementation of the storage interface that uses S3 as the store.
type S3Store struct {
	Bucket string
	Region string
}

// Get retrieves the data for the given key.
func (s *S3Store) Get(key string) ([]byte, error) {
	session, err := s.newSession()
	if err != nil {
		return nil, err
	}

	buffer := &aws.WriteAtBuffer{}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	if _, err := s3manager.NewDownloader(session).Download(buffer, input); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Put saves the data under the given key in S3.
func (s *S3Store) Put(key string, contentType string, data io.ReadSeeker) error {
	session, err := s.newSession()
	if err != nil {
		return err
	}

	input := &s3manager.UploadInput{
		Bucket:      aws.String(s.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
		Body:        data,
	}

	if _, err := s3manager.NewUploader(session).Upload(input); err != nil {
		return err
	}

	return nil
}

func (s *S3Store) newSession() (*session.Session, error) {
	config := &aws.Config{
		Region: aws.String(s.Region),
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
