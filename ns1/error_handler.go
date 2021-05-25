package ns1

import (
	ns1 "gopkg.in/ns1/ns1-go.v2/rest"
	"net/http"
)

// ConvertToNs1Error  convert messages that GoSDK client overrides to a verbose one
func ConvertToNs1Error(resp *http.Response, err error) error {
	if resp == nil {
		return err
	}

	if err == nil {
		return nil
	}

	if _, ok := err.(*ns1.Error); ok {
		return err
	}

	return &ns1.Error{Resp: resp, Message: err.Error()}
}
