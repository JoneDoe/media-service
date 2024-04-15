package cloud

import (
	"log"
	"mime"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"istorage/logger"
	"istorage/models"
	"istorage/s3"
)

func Store(attach *models.Attachment) error {
	file := attach.OriginalFile.Upload

	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	mimeType := mime.TypeByExtension(attach.OriginalFile.Ext())

	out, err := s3.NewUploader().Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s3.GetBucketName()),
		Key:         aws.String(attach.ToJson().FileSystemPath()),
		ACL:         aws.String("public-read"),
		Body:        src,
		ContentType: &mimeType,
	})

	if err != nil {
		return err
	}

	logger.Infof("Successfully uploaded %s", out.Location)

	return nil
}

func StoreFromLocal(file models.LocalFile) error {
	src, err := file.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	mimeType := file.MimeType()
	out, err := s3.NewUploader().Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s3.GetBucketName()),
		Key:         aws.String(file.FileSystemPath()),
		ACL:         aws.String("public-read"),
		Body:        src,
		ContentType: &mimeType,
	})

	if err != nil {
		return err
	}

	logger.Infof("Successfully uploaded %s", out.Location)

	return nil
}

func GetFullUrl(filePath string) string {
	return strings.Join([]string{s3.GetBucketUrl(), filePath}, "/")
}

func RemoveMedia(file *models.MediaFile) error {
	_, err := s3.New().DeleteObject(&awss3.DeleteObjectInput{
		Bucket: aws.String(s3.GetBucketName()),
		Key:    aws.String(file.FileSystemPath()),
	})

	if err != nil {
		return err
	}

	logger.Infof("Successfully deleted %s from %s", file.Name, s3.GetBucketName())

	return nil
}

func ReadFile(file *models.MediaFile) string {
	tmpFile, err := os.CreateTemp("", "s3-bucket.*")
	if err != nil {
		log.Fatal(err)
	}

	_, err = s3.NewDownloader().Download(tmpFile,
		&awss3.GetObjectInput{
			Bucket: aws.String(s3.GetBucketName()),
			Key:    aws.String(file.FileSystemPath()),
		})

	if err != nil {
		logger.Error(err)
	}

	return tmpFile.Name()
}

/*func Store(files []*multipart.FileHeader) {
	var wg sync.WaitGroup
	var m sync.Mutex

	filesLength := len(files)

	wg.Add(filesLength)

	objects := make([]s3manager.BatchUploadObject, 0)

	for _, file := range files {
		go func(file *multipart.FileHeader) {
			src, err := file.Open()
			defer src.Close()

			if err != nil {
				logger.Fatal(err)
			}

			m.Lock()

			objects = append(objects, s3manager.BatchUploadObject{
				Object: &s3manager.UploadInput{
					Bucket: aws.String(s3.GetBucketName()),
					Key:    aws.String(file.Filename),
					ACL:    aws.String("public-read"),
					Body:   src,
				},
			})

			m.Unlock()
			wg.Done()
		}(file)
	}

	wg.Wait()

	fmt.Println(objects)

	s3.BatchUploadedFiles(objects)
}*/
