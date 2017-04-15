package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"context"
	"io"
)

func NewClient() (Client, error) {
	conf, err := NewSettings()
	if err != nil {
		return Client{}, err
	}

	cred := credentials.NewStaticCredentials(conf.ID, conf.Secret, "")
	awsConf := aws.NewConfig().WithRegion(conf.Region).WithCredentials(cred)
	sess, err := session.NewSession(awsConf)
	if err != nil {
		return Client{}, err
	}

	return Client{
		conn:     s3.New(sess),
		settings: conf,
	}, nil
}

type Client struct {
	conn     *s3.S3
	settings Settings
}

func (c Client) Get(filename string) (io.ReadCloser, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.settings.ReadTimeout)
	out, err := c.conn.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(c.settings.Bucket),
			Key:    aws.String(filename),
		},
	)

	return out.Body, cancel, err
}
