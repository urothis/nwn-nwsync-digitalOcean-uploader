package main

import (
	"fmt"
	"os"

	"log"
	"path/filepath"
	"time"

	minio "github.com/minio/minio-go/v6"
	"github.com/schollz/progressbar/v2"
)

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
}

func fileCount(path string) (int, error) {
	i := 0
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {

			// return if directory
			if info, err := os.Stat(path); err == nil && info.IsDir() {
				return err
			}
			i++
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return i, nil
}

func main() {
	start := time.Now()

	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACCESS_KEY")
	secretAccessKey := os.Getenv("SECRET_KEY")

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, true)
	if err != nil {
		log.Fatal(err)
	}
	spaceName := os.Getenv("SPACE_NAME")

	folder := "/nwsync"
	totalCount, _ := fileCount(folder)
	log.Printf("%d files processing", int(totalCount))
	bar := progressbar.NewOptions(totalCount,
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowIts(),
	)
	err = filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			// return if directory
			if info, err := os.Stat(path); err == nil && info.IsDir() {
				return err
			}
			bar.Add(1)
			// upload object
			_, err = client.FPutObject(spaceName, "nwsync/"+os.Getenv("MODULENAME")+path[7:], path, minio.PutObjectOptions{ContentType: "application/gzip", UserMetadata: map[string]string{"x-amz-acl": "public-read"}})
			if err != nil {
				log.Fatalln(err)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	Totalelapsed := time.Since(start)
	bar.Finish()
	log.Printf("Successfully uploaded %d files | %s elapsed\n", int(count), fmtDuration(Totalelapsed))
}
