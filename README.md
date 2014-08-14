Example application for accessing the cloud storage from an App Engine application and also how to make it work in the App Engine SDK.

A note is that this code is a simple and very unsafe example, uploading it with access to your bucket basically grants full control to public.

Implemented example API's:
- Create/upload object
- List objects in bucket
- Delete object

Todo:
- Read object
- Get public link for object


## How to setup the project to work in App Engine SDK

1. Enable Cloud Integration in appengine console
2. Create Oauth 2.0 Client Id -> Service Account in Cloud Console.
3. Download the JSON and store it as `server/storage_dev_privatekey.json`
4. Edit the values `server/config.go` to match your project and bucket.