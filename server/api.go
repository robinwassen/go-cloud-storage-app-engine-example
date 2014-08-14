package server

import (
	"appengine"
	"appengine/blobstore"
	"appengine/image"
	"code.google.com/p/google-api-go-client/storage/v1beta2"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// List all objects in the Cloud Storage Bucket
func api_listObjects(w http.ResponseWriter, r *http.Request) {
	service, err := newStorageService(appengine.NewContext(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res, err := service.Objects.List(bucketName).Do(); err == nil {
		fmt.Fprintln(w, "Listing objects:")
		for _, object := range res.Items {
			fmt.Fprintln(w, object.Name)
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Get the image link for an object.
func api_imageLinkObject(w http.ResponseWriter, r *http.Request) {
	// We use the blob store + images to generate a serving url for the file
	c := appengine.NewContext(r)
	vars := mux.Vars(r)

	bk, err := blobstore.BlobKeyForFile(c, "/gs/"+bucketName+"/"+vars["key"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if imageUrl, err := image.ServingURL(c, bk, nil); err == nil {
		fmt.Fprint(w, imageUrl)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Read an object by key from the Cloud Storage Bucket
func api_readObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["key"])

	// TODO: Implement me!
}

// Create an object in the Cloud Storage Bucket
func api_createObject(w http.ResponseWriter, r *http.Request) {
	service, err := newStorageService(appengine.NewContext(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	object := &storage.Object{Name: fileHeader.Filename}

	if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Delete an object in the Cloud Storage Bucket
func api_deleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service, err := newStorageService(appengine.NewContext(r))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := service.Objects.Delete(bucketName, vars["key"]).Do(); err == nil {
		fmt.Fprintf(w, "Successfully deleted %s/%s.", bucketName, vars["key"])
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
