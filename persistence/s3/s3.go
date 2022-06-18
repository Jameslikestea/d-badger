package s3

//go:generate go run github.com/golang/mock/mockgen -destination ./mocks/s3.go -package mocks github.com/aws/aws-sdk-go/service/s3/s3iface S3API

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/persistence"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	_ persistence.Provider = &Service{}
	_ io.Writer            = &s3Writer{}
)

type Service struct {
	svc s3iface.S3API
	b   string
	p   string
}

func New(
	s *session.Session,
	bucket string,
	prefix string,
) *Service {
	svc := s3.New(s)

	return &Service{
		svc: svc,
		b:   bucket,
		p:   prefix,
	}
}

// Get implements persistence.Provider
func (s *Service) Get(key string) (io.Reader, error) {
	out, err := s.svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.b,
		Key:    aws.String(fmt.Sprintf("%s/%s", s.p, key)),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		if aerr.Code() != s3.ErrCodeNoSuchKey {
			log.Println(err)
			return nil, errors.KeyError
		}
		if aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, errors.NoSuchKey
		}
	}

	return out.Body, nil
}

// Put implements persistence.Provider
func (s *Service) Put(key string) (io.Writer, error) {
	k := fmt.Sprintf("%s/%s", s.p, key)

	wrtr := &s3Writer{
		s: s.svc,
		b: s.b,
		k: k,
	}

	return wrtr, nil
}

type s3Writer struct {
	s       s3iface.S3API
	b       string
	k       string
	content []byte
}

// Write implements io.Writer
func (s *s3Writer) Write(p []byte) (n int, err error) {
	s.content = append(s.content, p...)
	uploader := s3manager.NewUploaderWithClient(s.s)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      &s.b,
		Key:         &s.k,
		Body:        bytes.NewReader(s.content),
		ContentType: aws.String("binary/octet-stream"),
	})
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
