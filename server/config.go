package server

import (
	"code.google.com/p/google-api-go-client/storage/v1beta2"
)

const (
	// Change these variable to match your personal information.
	bucketName = "YOUR_BUCKET_NAME_TEST"
	projectID  = "YOUR_PROJECT_ID_TEST" // the app engine project id - used for creating buckets.
	scope      = storage.DevstorageFull_controlScope

	devStorageClientSecretPath = "storage_dev_clientsecret.json"
	devStorageSecretKeyPath    = "storage_dev_privatekey.pem"
)
