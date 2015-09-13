package deliver

import (
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
	"log"
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

		job.Result = Delivering(filename, bucketConfig)
		results <- job
	}
}

func Delivering(filename string, bucketConfig map[interface{}]interface{}) string {
	var (
		fullpath      string
		clearFilename string
		clearDir      string
	)

	if needConvering(filename) {
		a := strings.Split(filename, "/")
		a = append(a[:(len(a)-5)], a[(len(a)-1):]...)

		clearFilename = strings.Join(a, "/")

		clearDir = filepath.Dir(clearFilename)

		originalFilePath := "cache" + string(filepath.Separator) + clearFilename

		if !isCached(originalFilePath) {
			originalFile, err := getFromS3(bucketConfig, clearFilename)
			if err != nil {
				log.Printf("Delivering: %s: %s", filename, err)
				return fullpath
			} else {
				originalFilePath = createFile(originalFile, clearDir, clearFilename)
			}

		}

		fullpath, err := Convert(filename, originalFilePath)
		if err != nil {
			return ""
		} else {
			return fullpath
		}
	} else {
		dir := filepath.Dir(filename)

		file, err := getFromS3(bucketConfig, filename)
		if err != nil {
			log.Printf("Delivering: %s: %s", filename, err)
			return fullpath
		} else {
			fullpath = createFile(file, dir, filename)
		}
	}

	return fullpath
}

func getFromS3(bucketConfig map[interface{}]interface{}, filename string) (*[]byte, error) {
	auth := aws.Auth{
		AccessKey: bucketConfig["access_key_id"].(string),
		SecretKey: bucketConfig["secret_access_key"].(string),
	}
	euwest := aws.EUCentral

	connection := s3.New(auth, euwest)
	bucket, err := connection.Bucket(bucketConfig["bucket"].(string))
	if err != nil {
		return nil, err
	}

	file, err := bucket.Get(filename)

	if err != nil {
		return nil, err
	}
	return &file, nil
}

func needConvering(filename string) bool {
	a := strings.Split(filename, "/")

	if a[len(a)-3] == "s" && a[len(a)-5] == "gr" {
		return true
	} else {
		return false
	}
}
