package server

import (
	"code.google.com/p/google-api-go-client/storage/v1beta2"
)

const (
	// Change these variable to match your personal information.
	bucketName = "digiexam-play.appspot.com"
	projectID  = "digiexam-play" // the app engine project id - used for creating buckets.
	scope      = storage.DevstorageRead_writeScope

	devStorageKeyPath = "storage_dev_privatekey.json"
)
