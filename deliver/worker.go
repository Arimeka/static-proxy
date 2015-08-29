package deliver

import (
	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
	"log"
	"os"
	"path/filepath"
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
	var fullpath string = ""
	dir := filepath.Dir(filename)

	auth := aws.Auth{
		AccessKey: bucketConfig["access_key_id"].(string),
		SecretKey: bucketConfig["secret_access_key"].(string),
	}
	euwest := aws.EUCentral

	connection := s3.New(auth, euwest)
	bucket, err := connection.Bucket(bucketConfig["bucket"].(string))
	if err != nil {
		log.Println(err)
	} else {
		file, err := bucket.Get(filename)

		if err != nil {
			log.Println(err)
		} else {
			if dir != "." {
				os.MkdirAll("cache"+string(filepath.Separator)+dir, 0777)
			}

			fullpath = "cache" + string(filepath.Separator) + filename

			fi, err := os.Create(fullpath)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := fi.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			if _, err := fi.Write(file); err != nil {
				log.Fatal(err)
			}
		}
	}
	return fullpath
}
