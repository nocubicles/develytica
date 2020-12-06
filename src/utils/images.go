package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var bucketName = "skillbase.io"

func getImageKey(adID uint, fileName string) string {
	return "kuulutused/pildid/" + fmt.Sprint(adID) + "/" + fileName
}

func gets3Config() *aws.Config {

	var endPoint = os.Getenv("spacesEndPoint")
	var accessKey = os.Getenv("spacesAccessKey")
	var secret = os.Getenv("spacesSecret")
	return &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secret, ""),
		Endpoint:    aws.String(endPoint),
		Region:      aws.String("us-east-1"),
	}
}

func constructImageURL(objectKey string) string {
	return "https://" + bucketName + "." + os.Getenv("spacesEndPoint") + "/" + objectKey
}

func getClient() *s3.S3 {
	s3Session := session.New(gets3Config())

	s3Client := s3.New(s3Session)

	return s3Client
}

//GetAdImageUrls return all ad images that have been uploaded
func GetAdImageUrls(adID uint) []string {
	s3client := getClient()
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String("kuulutused/pildid/" + fmt.Sprint(adID)),
	}
	objects, err := s3client.ListObjects(input)
	if err != nil {
		fmt.Println(err)
	}

	adImageUrls := []string{}
	for _, object := range objects.Contents {
		imageURL := constructImageURL(fmt.Sprint(*object.Key))
		adImageUrls = append(adImageUrls, imageURL)
	}
	return adImageUrls
}

//UploadImage uploads image to bucket
func UploadImage(adID uint, file io.Reader, contentType string, fileName string, fileSize int64) {

	sess, err := session.NewSession(gets3Config())

	if err != nil {
		fmt.Println(err)
	}

	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.Concurrency = 1
	})

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(getImageKey(adID, fileName)),
		Body:   file,
		ACL:    aws.String("public-read"),
	}

	_, err = uploader.Upload(input)

	if err != nil {
		fmt.Println(err)
	}
}
