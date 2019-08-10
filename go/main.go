package main

import (
	"fmt"
	"os"

	"log"
	"path/filepath"
	"time"

	"github.com/janeczku/go-spinner"
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

	objectsCh := make(chan string)
	log.Println("Purging old nwsync data from cdn")
	s := spinner.StartNew("This may take some time...")
	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)

		doneCh := make(chan struct{})

		// Indicate to our routine to exit cleanly upon return.
		defer close(doneCh)

		// List all objects from a bucket-name with a matching prefix.
		for object := range client.ListObjects(spaceName, "nwsync/", true, doneCh) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object.Key
		}
	}()

	// Call RemoveObjects API
	errorCh := client.RemoveObjects(spaceName, objectsCh)

	// Print errors received from RemoveObjects API
	for e := range errorCh {
		log.Fatalln("Failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}
	s.Stop()
	log.Println("\nOld nwsync data purged")

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
	log.Printf("Successfully uploaded %d files | %s elapsed\n", int(totalCount), fmtDuration(Totalelapsed))
}
