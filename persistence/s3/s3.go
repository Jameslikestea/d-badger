package s3

import (
	"bytes"
	"fmt"
	"io"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/persistence"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	_ persistence.Provider = &Service{}
	_ io.Writer            = &s3Writer{}
)

type Service struct {
	s *session.Session
	b string
	p string
}

func New(
	s *session.Session,
	bucket string,
	prefix string,
) *Service {
	return &Service{
		s: s,
		b: bucket,
		p: prefix,
	}
}

// Get implements persistence.Provider
func (s *Service) Get(key string) (io.Reader, error) {
	svc := s3.New(s.s)

	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &s.b,
		Key:    aws.String(fmt.Sprintf("%s/%s", s.p, key)),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if !ok {
			return nil, err
		}
		if aerr.Code() != s3.ErrCodeNoSuchKey {
			return nil, errors.NoSuchKey
		}
	}

	return out.Body, nil
}

// Put implements persistence.Provider
func (s *Service) Put(key string) (io.Writer, error) {
	svc := s3.New(s.s)
	k := fmt.Sprintf("%s/%s", s.p, key)

	wrtr := &s3Writer{
		s: svc,
		b: s.b,
		k: k,
	}

	return wrtr, nil
}

type s3Writer struct {
	s *s3.S3
	b string
	k string
}

// Write implements io.Writer
func (s *s3Writer) Write(p []byte) (n int, err error) {
	_, err = s.s.PutObject(&s3.PutObjectInput{
		Bucket: &s.b,
		Key:    &s.k,
		Body:   bytes.NewReader(p),
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
