package server

import (
	"appengine"
	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/storage/v1beta2"
	"devserviceaccount"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()

	o := r.PathPrefix("/storage/object").Subrouter()
	o.HandleFunc("/", api_listObjects).Methods("GET")
	o.HandleFunc("/", api_createObject).Methods("POST")
	o.HandleFunc("/{key}/", api_readObject).Methods("GET")
	o.HandleFunc("/{key}/", api_deleteObject).Methods("DELETE")
	http.Handle("/", r)
}

// Creates a new Google Cloud Storage Client
func newStorageService(c appengine.Context) (*storage.Service, error) {
	var httpClient *http.Client
	var err error

	if appengine.IsDevAppServer() {
		httpClient, err = devserviceaccount.NewClient(c, devStorageKeyPath, scope)
	} else {
		httpClient, err = serviceaccount.NewClient(c, scope)
	}

	if err != nil {
		return nil, err
	}

	service, err := storage.New(httpClient)
	if err != nil {
		return nil, err
	}

	if err = initDefaultBucket(service); err != nil {
		return nil, err
	}

	return service, nil
}

// Creates the default bucket if it does not exist
func initDefaultBucket(service *storage.Service) error {
	if _, err := service.Buckets.Get(bucketName).Do(); err != nil {
		// No bucket found - Create a bucket.
		if _, err := service.Buckets.Insert(projectID, &storage.Bucket{Name: bucketName}).Do(); err != nil {
			return err
		}
	}

	return nil
}

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

// Read an object by key from the Cloud Storage Bucket
func api_readObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["key"])
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
		fmt.Fprintf(w, "Successfully deleted %s/%s during cleanup.\n\n", bucketName, vars["key"])
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
