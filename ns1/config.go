package ns1

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
)

type Config struct {
	Key       string
	Endpoint  string
	IgnoreSSL bool
}

// Client() returns a new NS1 client.
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

	// If NS1_DEBUG is set, log the requests and responses
	if os.Getenv("NS1_DEBUG") != "" {
		doer := ns1.Decorate(httpClient, Logging())
		client = ns1.NewClient(doer, decos...)
	} else {
		client = ns1.NewClient(httpClient, decos...)
	}

	client.RateLimitStrategySleep()

	log.Printf("[INFO] NS1 Client configured for Endpoint: %s", client.Endpoint.String())

	return client, nil
}

func Logging() ns1.Decorator {
	return func(d ns1.Doer) ns1.Doer {
		return ns1.DoerFunc(func(r *http.Request) (*http.Response, error) {
			log.Printf("%s: %s %s", r.UserAgent(), r.Method, r.URL)
			log.Printf("Headers: %s", r.Header)
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

func logRequest(original io.ReadCloser) (io.ReadCloser, error) {
	// Handle request contentType
	var bs bytes.Buffer
	defer original.Close()

	_, err := io.Copy(&bs, original)
	if err != nil {
		return nil, err
	}

	debugInfo := formatJSON(bs.Bytes())
	log.Printf("[DEBUG] Request Body: %s", debugInfo)

	return ioutil.NopCloser(strings.NewReader(bs.String())), nil

	return original, nil
}

func formatJSON(raw []byte) string {
	var rawData interface{}
	err := json.Unmarshal(raw, &rawData)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return string(raw)
	}
	pretty, err := json.MarshalIndent(rawData, "", "  ")
	if err != nil {
		log.Printf("[DEBUG] Unable to re-marshal JSON: %s", err)
		return string(raw)
	}

	return string(pretty)
}
