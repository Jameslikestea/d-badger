package s3

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/persistence/s3/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGet_Happy(t *testing.T) {
	type happyTest struct {
		name   string
		bucket string
		prefix string
		key    string
	}

	tests := []happyTest{
		{
			name:   "[HAPPY] Basic Run",
			bucket: "test-bucket",
			prefix: "test-prefix",
			key:    "something",
		},
		{
			name:   "[HAPPY] Deep Path",
			bucket: "test-bucket",
			prefix: "test-prefix",
			key:    "something/else.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockS3API(ctrl)

			s.EXPECT().GetObject(&s3.GetObjectInput{
				Bucket: &tt.bucket,
				Key:    aws.String(fmt.Sprintf("%s/%s", tt.prefix, tt.key)),
			}).Return(&s3.GetObjectOutput{
				Body: ioutil.NopCloser(strings.NewReader("content")),
			}, nil)

			svc := &Service{
				svc: s,
				b:   tt.bucket,
				p:   tt.prefix,
			}

			rdr, err := svc.Get(tt.key)
			assert.Nil(t, err)

			content, err := ioutil.ReadAll(rdr)
			assert.Nil(t, err)

			assert.Equal(t, "content", string(content))
		})
	}
}

func TestGet_Errors(t *testing.T) {
	type test struct {
		name      string
		awsErr    error
		expectErr error
	}

	var (
		bucket string = "test-bucket"
		prefix string = "test-prefix"
		key    string = "test-key"
	)

	var (
		s3nokey    = awserr.New(s3.ErrCodeNoSuchKey, s3.ErrCodeNoSuchKey, nil)
		s3nobucket = awserr.New(s3.ErrCodeNoSuchBucket, s3.ErrCodeNoSuchBucket, nil)
	)

	tests := []test{
		{
			name:      "[SAD] No Such Key",
			awsErr:    errors.NoSuchKey,
			expectErr: errors.NoSuchKey,
		},
		{
			name:      "[SAD] AWS Provides No Such Key",
			awsErr:    s3nokey,
			expectErr: errors.NoSuchKey,
		},
		{
			name:      "[SAD] AWS Provides No Such Bucket",
			awsErr:    s3nobucket,
			expectErr: errors.KeyError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mocks.NewMockS3API(ctrl)

			s.EXPECT().GetObject(&s3.GetObjectInput{
				Bucket: &bucket,
				Key:    aws.String(fmt.Sprintf("%s/%s", prefix, key)),
			}).Return(&s3.GetObjectOutput{}, tt.awsErr)

			svc := &Service{
				svc: s,
				b:   bucket,
				p:   prefix,
			}

			_, err := svc.Get(key)
			assert.Equal(t, tt.expectErr, err)
		})
	}
}
