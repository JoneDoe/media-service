package s3

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	config Config
	sess   *session.Session
)

type Config struct {
	S3host, S3bucket string
}

func SetConfig(cfg Config) {
	config = cfg
}

func GetBucketName() string {
	return config.S3bucket
}

func GetBucketUrl() string {
	return strings.Join([]string{endpoints.AddScheme(config.S3host, false), config.S3bucket}, "/")
}

func New() *s3.S3 {
	initSession()

	return s3.New(sess)
}

func NewDownloader() *s3manager.Downloader {
	initSession()

	return s3manager.NewDownloader(sess)
}

func NewUploader() *s3manager.Uploader {
	initSession()

	return s3manager.NewUploader(sess)
}

func initSession() {
	cfg := &aws.Config{
		S3ForcePathStyle:    aws.Bool(true),
		LowerCaseHeaderMaps: aws.Bool(true),
		Region:              aws.String(endpoints.UsEast1RegionID),
		Credentials:         credentials.NewEnvCredentials(),
	}

	if config.S3host != "" {
		cfg.Endpoint = aws.String(endpoints.AddScheme(config.S3host, false))
	}

	sess = session.Must(session.NewSession(cfg))
}
