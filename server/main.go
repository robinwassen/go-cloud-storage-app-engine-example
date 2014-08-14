package server

import (
	"appengine"
	"code.google.com/p/goauth2/appengine/serviceaccount"
	"code.google.com/p/google-api-go-client/storage/v1beta2"
	"devserviceaccount"
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
