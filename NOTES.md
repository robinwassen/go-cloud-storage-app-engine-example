## Service Account Authentication

Service accounts are special accounts that represent software rather than people. 

https://developers.google.com/storage/docs/authentication#service_accounts
https://code.google.com/p/goauth2/source/browse/appengine/serviceaccount/serviceaccount.go?r=ecc4c1308422bb3ab05a022036dfdf05f2c05272
https://groups.google.com/forum/#!topic/google-appengine-go/qZDQHkzEFMU


1. Enable Cloud Integration in appengine console
2. Create Oauth 2.0 Client Id -> Service Account in Cloud Console.
3. Download the JSON and store it as `server/storage_dev_privatekey.json`
