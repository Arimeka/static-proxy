package deliver

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WorkerPool(n int) (jobs chan *Job, results chan *Job) {
	jobs = make(chan *Job)
	results = make(chan *Job)

	for i := 0; i < n; i++ {
		go Worker(jobs, results)
	}

	return
}

func Worker(jobs chan *Job, results chan *Job) {
	for job := range jobs {
		filename := job.Filename
		bucketConfig := job.BucketConfig.(map[interface{}]interface{})

		defaults.DefaultConfig.Credentials = credentials.NewStaticCredentials(bucketConfig["access_key_id"].(string), bucketConfig["secret_access_key"].(string), "")
		defaults.DefaultConfig.Region = aws.String("eu-central-1")

		job.Result = Delivering(filename, bucketConfig)
		results <- job
	}
}

func Delivering(filename string, bucketConfig map[interface{}]interface{}) string {
	var (
		fullpath      string
		clearFilename string
	)

	if needConvering(filename) {
		a := strings.Split(filename, "/")
		a = append(a[:(len(a)-5)], a[(len(a)-1):]...)

		clearFilename = strings.Join(a, "/")

		originalFilePath := "cache" + string(filepath.Separator) + clearFilename

		if !isCached(originalFilePath) {
			_, err := getFromS3(bucketConfig, clearFilename)
			if err != nil {
				log.Printf("Delivering: %s: %s", filename, err)
				return fullpath
			}
		}

		fullpath, err := Convert(filename, originalFilePath)
		if err != nil {
			return ""
		} else {
			return fullpath
		}
	} else {
		fullpath, err := getFromS3(bucketConfig, filename)
		if err != nil {
			log.Printf("Delivering: %s: %s", filename, err)
			return fullpath
		}
	}

	return fullpath
}

func getFromS3(bucketConfig map[interface{}]interface{}, filename string) (string, error) {
	svc := s3.New(nil)
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketConfig["bucket"].(string)),
		Key:    aws.String(filename),
	})
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(filename)

	if dir != "." {
		os.MkdirAll("cache"+string(filepath.Separator)+dir, 0755)
	}

	fullpath := "cache" + string(filepath.Separator) + filename

	file, err := os.Create(fullpath)
	if err != nil {
		log.Fatal("Failed to create file", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if _, err := io.Copy(file, result.Body); err != nil {
		log.Fatal("Failed to copy object to file", err)
	}
	result.Body.Close()

	return fullpath, nil
}

func needConvering(filename string) bool {
	a := strings.Split(filename, "/")

	if a[len(a)-3] == "s" && a[len(a)-5] == "gr" {
		return true
	} else {
		return false
	}
}

func isCached(cachedPath string) bool {
	if _, err := os.Stat(cachedPath); err == nil {
		return true
	} else {
		return false
	}
}
