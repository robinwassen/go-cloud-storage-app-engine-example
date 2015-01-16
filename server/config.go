package server

import (
	"code.google.com/p/google-api-go-client/storage/v1beta2"
)

const (
	// Change these variable to match your personal information.
	bucketName        = "digiexam-play.appspot.com"   // The bucket name for Cloud Integrations is the same as the appspot url
	projectID         = "digiexam-play"               // Project ID is the one found in the cloud console
	devStorageKeyPath = "storage_dev_privatekey.json" // This file is downloaded from the cloud console

	// Set this according to: 
	// https://github.com/google/google-api-go-client/blob/master/storage/v1/storage-gen.go
	// write_scope will work for all file uploading
	// use full_control if you wish to pass along metadata in the attributes object e.g.
	//
	// object := &storage.Object{
	//	Name:        fileHeader.Filename,
	//	ContentType: YOUR_CONTENT_TYPE,
	//}

	scope = storage.DevstorageRead_writeScope // TODO: Figure out exactly what this does.
)
