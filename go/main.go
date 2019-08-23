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

// Client Connection
type Client struct {
	client *minio.Client
	Bucket string
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
}

func localFileCount(path string) (int, error) {
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

func purgeOldNWSyncData(client *minio.Client, spaceName, objectPrefix string) int {
	log.Println("Purging old nwsync data from cdn")

	doneCh := make(chan struct{})
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	objectsCh := make(chan string)
	count := 0
	for object := range client.ListObjects(spaceName, objectPrefix, true, doneCh) {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		count++
		objectsCh <- object.Key
	}

	// Call RemoveObjects API
	errorCh := client.RemoveObjects(spaceName, objectsCh)
	// Print errors received from RemoveObjects API
	for e := range errorCh {
		log.Fatalln("Failed to remove " + e.ObjectName + ", error: " + e.Err.Error())
	}

	log.Println("\n Old nwsync data purged")

	return count
}

func uploadNewNWSyncData(client *minio.Client, folder, spaceName string) (int, error) {
	totalCount, err := localFileCount("/nwsync")
	if err != nil {
		return 0, err
	}
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
	bar.Finish()
	return totalCount, err
}

func main() {
	start := time.Now()

	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACCESS_KEY")
	secretAccessKey := os.Getenv("SECRET_KEY")
	spaceName := os.Getenv("SPACE_NAME")
	localPath := os.Getenv("NWSYNCFILEPATH")
	objectPrefix := "nwsync/" + os.Getenv("MODULENAME") + "/"

	// Initiate a client
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, true)
	if err != nil {
		log.Fatal(err)
	}

	// clean the old nwsync data
	deletedCount := purgeOldNWSyncData(client, spaceName, objectPrefix)
	if err != nil {
		log.Fatal(err)
	}

	// upload new files
	totalCount, err := uploadNewNWSyncData(client, localPath, spaceName)
	if err != nil {
		log.Fatal(err)
	}

	Totalelapsed := time.Since(start)
	log.Printf("Successfully deleted %d old files | uploaded %d new files | %s elapsed\n", int(deletedCount), int(totalCount), fmtDuration(Totalelapsed))
}
