// +build appengine

// Holger Knauer
// https://groups.google.com/forum/#!msg/google-appengine-go/HFT3wuYG0Jg/0yArpTHspD4J

package devserviceaccount

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/goauth2/oauth/jwt"
)

// NewClient returns an *http.Client authorized for the
// given scopes with the service account define by the given pem key file and client-secret json file
// Tokens are cached in memcache until they expire.
func NewClient(c appengine.Context, clientSecretFileName string, scopes ...string) (*http.Client, error) {
	t := &transport{
		Context: c,
		Scopes:  scopes,
		Transport: &urlfetch.Transport{
			Context:                       c,
			Deadline:                      0,
			AllowInvalidServerCertificate: false,
		},
		ClientSecretFileName: clientSecretFileName,
		TokenCache: &cache{
			Context: c,
			Key:     "goauth2_serviceaccount_" + strings.Join(scopes, "_"),
		},
	}
	// Get the initial access token.
	if err := t.FetchToken(); err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: t,
	}, nil
}

// transport is an oauth.Transport with a custom Refresh and RoundTrip implementation.
type transport struct {
	*oauth.Token

	Context    appengine.Context
	Scopes     []string
	Transport  http.RoundTripper
	TokenCache oauth.Cache

	ClientSecretFileName string
}

type keyConfig struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	Type         string `json:"type"`
}

func (t *transport) Refresh() error {
	// Before we can get a new oauth token, we have to produce a jwt token
	//
	// the oAuth token returned does not contain a refresh token and expires after one hour,
	// which is the same expiration the JWT token provides. Therefore it does not make sense
	// to keep the JWT token around, just rebuild it on every refresh

	// Read the secret file bytes into the config.
	secretFileBytes, err := ioutil.ReadFile(t.ClientSecretFileName)
	if err != nil {
		return fmt.Errorf("Error reading file %q: %v", t.ClientSecretFileName, err)
	}

	var config keyConfig

	err = json.Unmarshal(secretFileBytes, &config)
	if err != nil {
		return fmt.Errorf("Failed to unmarshall json in %q: %v", t.ClientSecretFileName, err)
	}

	// Craft the ClaimSet and JWT token.
	jwtToken := jwt.NewToken(config.ClientEmail, strings.Join(t.Scopes, " "), []byte(config.PrivateKey))
	//jwtToken.ClaimSet.Aud = config.Web.TokenURI // just in case: assume that TokenURI from secret.json is more current than the default jwt package's

	// assert the jwtToken to get the oAuth Token
	client := urlfetch.Client(t.Context)
	t.Token, err = jwtToken.Assert(client)
	if err != nil {
		return fmt.Errorf("Assert jwt token error: %v", err)
	}

	if t.TokenCache != nil {
		// Cache the token and ignore error (as we can always get a new one).
		t.TokenCache.PutToken(t.Token)
	}
	return nil
}

// Fetch token from cache or generate a new one if cache miss or expired.
func (t *transport) FetchToken() error {
	// Try to get the Token from the cache if enabled.
	if t.Token == nil && t.TokenCache != nil {
		// Ignore cache error as we can always get a new token with Refresh.
		t.Token, _ = t.TokenCache.Token()
	}

	// Get a new token using Refresh in case of a cache miss of if it has expired.
	if t.Token == nil || t.Expired() {
		if err := t.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

// RoundTrip issues an authorized HTTP request and returns its response.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.FetchToken(); err != nil {
		return nil, err
	}

	// To set the Authorization header, we must make a copy of the Request
	// so that we don't modify the Request we were given.
	// This is required by the specification of http.RoundTripper.
	newReq := cloneRequest(req)
	newReq.Header.Set("Authorization", "Bearer "+t.AccessToken)

	// Make the HTTP request.
	return t.Transport.RoundTrip(newReq)
}

// cache implementss TokenCache using memcache to store AccessToken
// for the application service account.
type cache struct {
	Context appengine.Context
	Key     string
}

func (m cache) Token() (*oauth.Token, error) {

	tok := new(oauth.Token)
	_, err := memcache.Gob.Get(m.Context, m.Key, tok)
	if err != nil {
		return nil, err
	}

	return tok, nil
}

func (m cache) PutToken(tok *oauth.Token) error {
	return memcache.Gob.Set(m.Context, &memcache.Item{
		Key: m.Key,
		Object: oauth.Token{
			AccessToken: tok.AccessToken,
			Expiry:      tok.Expiry,
		},
		Expiration: tok.Expiry.Sub(time.Now()),
	})
}
