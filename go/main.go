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

// pretty prints the elapsed time
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh%02dm", h, m)
}

// count the nwsync files
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

// delete old nwsync data from cdn.
func purgeOldNWSyncData(client *minio.Client, spaceName, objectPrefix string) int {
	log.Println("Purging old nwsync data from cdn")

	doneCh := make(chan struct{})
	defer close(doneCh)

	// delete path
	path := os.Getenv("MODULE_NAME") + "/";

	// List all objects from a bucket-name with a matching prefix.
	count := 0
	for object := range client.ListObjects(spaceName, objectPrefix, true, doneCh) {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}
		client.RemoveObject(spaceName, path + object.Key)
		count++
	}

	return count
}

func uploadNewNWSyncData(client *minio.Client, folder, spaceName string) (int, error) {
	// count files
	totalCount, err := localFileCount(folder) // "/nwsync"
	if err != nil {
		return 0, err
	}
	log.Printf("%d files processing", int(totalCount))

	// progress bar setup
	bar := progressbar.NewOptions(totalCount,
		   progressbar.OptionSetRenderBlankState(true),
		   progressbar.OptionShowIts(),
	)

	// walk the filepath
	err = filepath.Walk(folder,
		func(path string, info os.FileInfo, err error) error {
			// return if directory
			if info, err := os.Stat(path); err == nil && info.IsDir() {
				return err
			}
			
			// add increment to progress bar
			bar.Add(1)

			// upload object
			_, err = client.FPutObject(spaceName, os.Getenv("MODULE_NAME") + "/" + path, path, minio.PutObjectOptions{ContentType: "application/gzip", UserMetadata: map[string]string{"x-amz-acl": "public-read"}})
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
	// grab the current time of start
	start := time.Now()

	// endpoint URL to object storage service
	endpoint := os.Getenv("ENDPOINT")
	// access key is the user ID that uniquely identifies your account.
	accessKeyID := os.Getenv("ACCESS_KEY")
	// secret key is the password to your account.
	secretAccessKey := os.Getenv("SECRET_KEY")

	// name of space
	spaceName := os.Getenv("SPACE_NAME")

	// prefix for multiple server nwsync in one container
	objectPrefix := os.Getenv("MODULE_NAME") + "/nwsync/"

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
	uploadedCount, err := uploadNewNWSyncData(client, "nwsync", spaceName)
	if err != nil {
		log.Fatal(err)
	}

	Totalelapsed := time.Since(start)
	log.Printf("Successfully deleted %d old nwsync files | uploaded %d new nwsync files | %s elapsed\n", int(deletedCount), int(uploadedCount), fmtDuration(Totalelapsed)) 
}
