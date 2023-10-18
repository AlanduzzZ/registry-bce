package tencentcdn

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	dcontext "github.com/distribution/distribution/v3/context"
	storagedriver "github.com/distribution/distribution/v3/registry/storage/driver"
	storagemiddleware "github.com/distribution/distribution/v3/registry/storage/driver/middleware"

	"tencentyun/auth"
)

// aliCDNStorageMiddleware provides a simple implementation of layerHandler that
// constructs temporary signed AliCDN URLs from the storagedriver layer URL,
// then issues HTTP Temporary Redirects to this AliCDN content URL.
type bceCDNStorageMiddleware struct {
	storagedriver.StorageDriver
	baseURL   string
	urlSigner *auth.URLSigner
	duration  time.Duration
}

var _ storagedriver.StorageDriver = &bceCDNStorageMiddleware{}

// newAliCDNStorageMiddleware constructs and returns a new AliCDN
// StorageDriver implementation.
// Required options: baseurl, authtype, privatekey
// Optional options: duration
func newBceCDNStorageMiddleware(storageDriver storagedriver.StorageDriver, options map[string]interface{}) (storagedriver.StorageDriver, error) {
	// parse baseurl
	base, ok := options["baseurl"]
	if !ok {
		return nil, fmt.Errorf("no baseurl provided")
	}
	baseURL, ok := base.(string)
	if !ok {
		return nil, fmt.Errorf("baseurl must be a string")
	}
	if !strings.Contains(baseURL, "://") {
		baseURL = "https://" + baseURL
	}
	if _, err := url.Parse(baseURL); err != nil {
		return nil, fmt.Errorf("invalid baseurl: %v", err)
	}

	// parse authtype
	at, ok := options["authtype"]
	if !ok {
		return nil, fmt.Errorf("no authtype provided")
	}
	authType, ok := at.(string)
	if !ok {
		return nil, fmt.Errorf("authtype must be a string")
	}
	if authType != "a" && authType != "b" && authType != "c" {
		return nil, fmt.Errorf("invalid authentication type")
	}

	// parse privatekey
	pk, ok := options["privatekey"]
	if !ok {
		return nil, fmt.Errorf("no privatekey provided")
	}
	privateKey, ok := pk.(string)
	if !ok {
		return nil, fmt.Errorf("privatekey must be a string")
	}

	urlSigner := auth.NewURLSigner(authType, privateKey)

	// parse duration
	duration := 20 * time.Minute
	d, ok := options["duration"]
	if ok {
		switch d := d.(type) {
		case time.Duration:
			duration = d
		case string:
			dur, err := time.ParseDuration(d)
			if err != nil {
				return nil, fmt.Errorf("invalid duration: %s", err)
			}
			duration = dur
		}
	}

	return &bceCDNStorageMiddleware{
		StorageDriver: storageDriver,
		baseURL:       baseURL,
		urlSigner:     urlSigner,
		duration:      duration,
	}, nil
}

// URLFor attempts to find a url which may be used to retrieve the file at the given path.
func (bc *bceCDNStorageMiddleware) URLFor(ctx context.Context, path string, options map[string]interface{}) (string, error) {
	if !(bc.StorageDriver.Name() != "s3" || bc.StorageDriver.Name() != "cos") {
		dcontext.GetLogger(ctx).Warn("the TencentCloud middleware does not support this backend storage driver")
		return bc.StorageDriver.URLFor(ctx, path, options)
	}
	bcURL, err := bc.urlSigner.Sign(bc.baseURL+path, time.Now().Add(bc.duration))
	dcontext.GetLogger(ctx).Info(bcURL)
	if err != nil {
		return "", err
	}
	return bcURL, nil
}

// init registers the alicdn layerHandler backend.
func init() {
	storagemiddleware.Register("tencentcdn", newBceCDNStorageMiddleware)
}
