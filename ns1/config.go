package ns1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

var (
	clientVersion     = "1.12.7"
	providerUserAgent = "terraform-ns1" + "/" + clientVersion
)

// Config for NS1 API
type Config struct {
	Key                  string
	Endpoint             string
	IgnoreSSL            bool
	EnableDDI            bool
	RateLimitParallelism int
}

// Client returns a new NS1 client.
func (c *Config) Client() (*ns1.Client, error) {
	var client *ns1.Client
	httpClient := &http.Client{}
	decos := []func(*ns1.Client){}

	if c.Key == "" {
		return nil, errors.New(`no valid credential sources found for NS1 Provider.
  Please see https://terraform.io/docs/providers/ns1/index.html for more information on
  providing credentials for the NS1 Provider`)
	}

	decos = append(decos, ns1.SetAPIKey(c.Key))
	if c.Endpoint != "" {
		decos = append(decos, ns1.SetEndpoint(c.Endpoint))
	}
	if c.IgnoreSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient.Transport = tr
	}

	if c.EnableDDI {
		decos = append(decos, ns1.SetDDIAPI())
	}

	// If NS1_DEBUG is set, define custom Doer to log HTTP requests made by SDK
	if os.Getenv("NS1_DEBUG") != "" {
		doer := ns1.Decorate(httpClient, Logging())
		client = ns1.NewClient(doer, decos...)
	} else {
		client = ns1.NewClient(httpClient, decos...)
	}

	if parallelism := c.RateLimitParallelism; parallelism > 0 {
		client.RateLimitStrategyConcurrent(parallelism)
	} else {
		client.RateLimitStrategySleep()
	}

	UA := providerUserAgent + "_" + client.UserAgent
	log.Printf("[INFO] NS1 Client configured for Endpoint: %s, versions %s", client.Endpoint.String(), UA)
	if localUA := os.Getenv("NS1_TF_USER_AGENT"); localUA != "" {
		client.UserAgent = localUA
	} else {
		client.UserAgent = UA
	}

	return client, nil
}

// Logging returns a ns1.Decorator with a ns1.Doer lambda that logs HTTP requests
func Logging() ns1.Decorator {
	return func(d ns1.Doer) ns1.Doer {
		return ns1.DoerFunc(func(r *http.Request) (*http.Response, error) {
			log.Printf("[DEBUG] %s: %s %s", r.UserAgent(), r.Method, r.URL)
			heads := r.Header.Clone()
			heads["X-Nsone-Key"] = []string{"<redacted>"}
			log.Printf("[DEBUG] Headers: %s", heads)
			var err error
			if r.Body != nil {
				r.Body, err = logRequest(r.Body)
				if err != nil {
					return nil, err
				}
			}
			return d.Do(r)
		})
	}
}

// logRequest logs a HTTP request and returns a copy that can be read again
func logRequest(original io.ReadCloser) (io.ReadCloser, error) {
	// Handle request contentType
	var bs bytes.Buffer
	defer original.Close()

	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	debugInfo, err := formatJSON(bs.Bytes())
	if err == nil {
		log.Printf("[DEBUG] Request Body: %s", debugInfo)
	}

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil
}

// formatJSON attempts to format a byte slice as indented JSON for pretty printing
func formatJSON(raw []byte) (string, error) {
	var rawData interface{}
	err := json.Unmarshal(raw, &rawData)
	if err != nil {
		return string(raw), fmt.Errorf("unable to parse JSON: %s", err)
	}
	pretty, err := json.MarshalIndent(rawData, "", "  ")
	if err != nil {
		return string(raw), fmt.Errorf("unable to re-marshal JSON: %s", err)
	}

	return string(pretty), nil
}
