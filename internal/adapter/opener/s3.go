package opener

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/cm-dev/template/internal/ports"
)

type OpenerS3 struct {
	client     *s3.Client
	keyPrefix  string
	bucketName string
}

var _ ports.Opener = (*OpenerS3)(nil)

func NewOpenerS3(bucketName, keyPrefix string) (*OpenerS3, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-west-2"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &OpenerS3{
		client,
		keyPrefix,
		bucketName,
	}, nil
}

func (o *OpenerS3) Open(ctx context.Context, name string) (io.ReadCloser, error) {
	listObjectsInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(o.bucketName),
		Prefix: aws.String(path.Join(o.keyPrefix, name)),
	}
	res, err := o.client.ListObjectsV2(ctx, listObjectsInput)
	if err != nil {
		return nil, fmt.Errorf("unable to list items in bucket %q: %v", o.bucketName, err)
	}

	var readers []io.ReadCloser
	for _, item := range res.Contents {
		getObjectInput := &s3.GetObjectInput{
			Bucket: aws.String(o.bucketName),
			Key:    item.Key,
		}
		obj, err := o.client.GetObject(ctx, getObjectInput)
		if err != nil {
			return nil, fmt.Errorf("unable to get object %q from bucket %q: %v", *item.Key, o.bucketName, err)
		}

		readers = append(readers, obj.Body)
	}

	return MultiReadCloser(readers...), nil
}
