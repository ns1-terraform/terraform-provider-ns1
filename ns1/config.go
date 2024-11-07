package ns1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

var (
	clientVersion     = "2.5.0"
	providerUserAgent = "tf-ns1" + "/" + clientVersion
	defaultRetryMax   = 3
)

// Config for NS1 API
type Config struct {
	Key                  string
	Endpoint             string
	IgnoreSSL            bool
	RateLimitParallelism int
	RetryMax             int
	UserAgent            string
}

// Client returns a new NS1 client.
func (c *Config) Client() (*ns1.Client, error) {
	var client *ns1.Client

	httpClient := &http.Client{}
	retryMax := c.RetryMax
	if retryMax >= 0 {
		if retryMax == 0 {
			retryMax = defaultRetryMax
		}
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = retryMax
		retryClient.Logger = nil
		httpClient = retryClient.StandardClient()
	}
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
	if len(c.UserAgent) > 0 {
		client.UserAgent = c.UserAgent
	} else {
		client.UserAgent = UA
	}
	log.Printf("[INFO] NS1 Client configuration: endpoint: %s, version %s, retries %d, User-Agent %s", client.Endpoint.String(), clientVersion, c.RetryMax, client.UserAgent)

	return client, nil
}

// Logging returns a ns1.Decorator with a ns1.Doer lambda that logs HTTP requests
func Logging() ns1.Decorator {
	return func(d ns1.Doer) ns1.Doer {
		return ns1.DoerFunc(func(r *http.Request) (*http.Response, error) {
			msgs := []string{}
			msgs = append(msgs, fmt.Sprintf("[DEBUG] HTTP %s %s", r.Method, r.URL))
			heads := r.Header.Clone()
			heads["X-Nsone-Key"] = []string{"<redacted>"}
			msgs = append(msgs, fmt.Sprintf("[DEBUG] HTTP Headers: %s", heads))
			var err error
			if r.Body != nil {
				var bodymsg string
				r.Body, bodymsg, err = logRequest(r.Body)
				if err != nil {
					return nil, err
				}
				msgs = append(msgs, bodymsg)
			}
			requestTime := time.Now()
			response, rerr := d.Do(r)
			responseTime := time.Now()
			dump := []byte{}
			if response != nil {
				dump, _ = httputil.DumpResponse(response, true)
			}
			for _, m := range msgs {
				log.Print(m)
			}
			log.Printf("[DEBUG] HTTP Response (requested at %s, received at %s): %s", requestTime.Format(time.StampMilli), responseTime.Format(time.StampMilli), dump)
			return response, rerr
		})
	}
}

// logRequest logs a HTTP request and returns a copy that can be read again
func logRequest(original io.ReadCloser) (io.ReadCloser, string, error) {
	// Handle request contentType
	var bs bytes.Buffer
	defer original.Close()

	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, "", err
	}

	msg := ""
	debugInfo, err := formatJSON(bs.Bytes())
	if err == nil {
		msg = fmt.Sprintf("[DEBUG] HTTP Request Body: %s", debugInfo)
	}

	return io.NopCloser(strings.NewReader(bs.String())), msg, nil
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
