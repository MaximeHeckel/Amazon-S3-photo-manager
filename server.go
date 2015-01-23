package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
)

var (
	bucketName string
	fileName   string
	imageFile  string
)

func init() {
	flag.StringVar(&bucketName, "b", "", "Bucket Name")
	flag.StringVar(&imageFile, "i", "", "Image File")
}

func readFile(file string) ([]byte, string) {
	imgfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer imgfile.Close()

	fileInfo, _ := imgfile.Stat()
	filename := fileInfo.Name()
	size := fileInfo.Size()
	bytes := make([]byte, size)

	buffer := bufio.NewReader(imgfile)
	_, err = buffer.Read(bytes)

	return bytes, filename
}

func main() {
	flag.Parse()

	//AWS AUTH
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}

	bytes, filename := readFile(imageFile)
	//Set new S3 connection
	s := s3.New(auth, aws.USEast)

	//Set connection to bucket
	bucket := s.Bucket(bucketName)

	//File Info
	filetype := http.DetectContentType(bytes)
	err = bucket.Put(filename, bytes, filetype, s3.BucketOwnerFull)
	if err != nil {
		panic(err.Error())
	}
}
