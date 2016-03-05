package deliver

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"static-proxy/settings"

	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
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

		if settings.Config.S3Config.Hosts[job.Host] == nil || settings.Config.ValidSizes.Sizes[job.Host] == nil {
			job.Result = ""
			results <- job
			continue
		}

		bucketConfig := settings.Config.S3Config.Hosts[job.Host]
		validSizes := settings.Config.ValidSizes.Sizes[job.Host]

		job.Result = Delivering(filename, bucketConfig, validSizes)
		results <- job
	}
}

func Delivering(filename string, bucketConfig map[string]string, validSizes []string) (fullpath string) {
	var (
		clearFilename    string
		originalFilePath string
		err              error
	)

	if needConvering(filename) {
		a := strings.Split(filename, "/")

		if sizeValid(a[len(a)-2], validSizes) == false {
			log.Printf("Delivering: %s: %s: %s", filename, "invalid size", a[len(a)-2])
			return ""
		}

		if len(a) >= 5 && a[len(a)-5] == "gr" {
			a = append(a[:(len(a)-5)], a[(len(a)-1):]...)
		} else {
			a = append(a[:(len(a)-3)], a[(len(a)-1):]...)
		}

		clearFilename = strings.Join(a, "/")

		originalFilePath = "cache" + string(filepath.Separator) + clearFilename

		if !isCached(originalFilePath) {
			_, err = getFromS3(bucketConfig, clearFilename)
			if err != nil {
				log.Printf("Delivering: %s: %s", filename, err)
				return
			}
		}

		fullpath, err = Convert(filename, originalFilePath)
		if err != nil {
			log.Printf("Converting: %s: %s", filename, err)
			return ""
		}
		return
	}
	fullpath, err = getFromS3(bucketConfig, filename)
	if err != nil {
		log.Printf("Delivering: %s: %s", filename, err)
		return ""
	}

	return
}

func getFromS3(bucketConfig map[string]string, filename string) (fullpath string, err error) {
	cred := credentials.NewStaticCredentials(bucketConfig["access_key_id"], bucketConfig["secret_access_key"], "")

	svc := s3.New(session.New(), &aws.Config{Credentials: cred, Region: aws.String("eu-central-1")})
	result, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketConfig["bucket"]),
		Key:    aws.String(filename),
	})
	if err != nil {
		return
	}

	dir := filepath.Dir(filename)

	if dir != "." {
		os.MkdirAll("cache"+string(filepath.Separator)+dir, 0755)
	}

	fullpath = "cache" + string(filepath.Separator) + filename

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

	return
}

func needConvering(filename string) bool {
	a := strings.Split(filename, "/")

	if len(a) >= 5 && a[len(a)-3] == "s" && a[len(a)-5] == "gr" {
		return true
	}
	if len(a) >= 3 && a[len(a)-3] == "s" {
		return true
	}

	return false
}

func sizeValid(size string, sizes []string) bool {
	sort.Strings(sizes)
	i := sort.Search(len(sizes), func(i int) bool { return sizes[i] >= size })
	if i < len(sizes) && sizes[i] == size {
		return true
	}
	return false
}

func isCached(cachedPath string) bool {
	if _, err := os.Stat(cachedPath); err == nil {
		return true
	}
	return false
}
