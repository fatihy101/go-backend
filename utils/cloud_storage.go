package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"enstrurent.com/server/flags"
)

func DeleteFromCloud(object string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(flags.GetConfig().GCP_BUCKET_NAME)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := bucket.Object(object)
	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", object, err)
	}
	return nil
}

// uploadFile uploads an object.
func UploadToCloud(r *http.Request) ([]string, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(flags.GetConfig().GCP_BUCKET_NAME)

	// Get Image from form.
	err = r.ParseMultipartForm(200000)
	if err != nil {
		return nil, err
	}

	formdata := r.MultipartForm
	files := formdata.File["images"] // grab the filenames
	var imageNames []string

	for _, fileHeader := range files {
		// Open the file.
		f, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if fileHeader.Size == 0 {
			return nil, errors.New("there's no photo or corrupted file")
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		fileName := fmt.Sprintf("p-%v.%v", rand.Int(), extractExtension(fileHeader.Filename))
		wc := bucket.Object(fileName).NewWriter(ctx)

		if _, err = io.Copy(wc, f); err != nil {
			return nil, fmt.Errorf("io.Copy: %v", err)
		}

		if err := wc.Close(); err != nil {
			return nil, fmt.Errorf("Writer.Close: %v", err)
		}

		imageNames = append(imageNames, fileName)
	}

	return imageNames, nil
}

func DownloadFromCloud(object string) ([]byte, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(flags.GetConfig().GCP_BUCKET_NAME)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	rc, err := bucket.Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

func extractExtension(val string) string {
	splitted := strings.Split(val, ".")
	return splitted[len(splitted)-1]
}
