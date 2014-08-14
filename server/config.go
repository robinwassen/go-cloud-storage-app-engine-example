package server

import (
	"code.google.com/p/google-api-go-client/storage/v1beta2"
)

const (
	// Change these variable to match your personal information.
	bucketName        = "digiexam-play.appspot.com"   // The bucket name for Cloud Integrations is the same as the appspot url
	projectID         = "digiexam-play"               // Project ID is the one found in the cloud console
	devStorageKeyPath = "storage_dev_privatekey.json" // This file is downloaded from the cloud console

	// No touching!
	scope = storage.DevstorageRead_writeScope // TODO: Figure out exactly what this does.
)
