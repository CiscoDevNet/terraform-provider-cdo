// request encapsulate the request to be sent
// it can be replayed by calling Send() multiple times
// it is not safe for concurrent use
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/jsonutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

type Request struct {
	config     cdo.Config
	httpClient *http.Client
	logger     *log.Logger

	ctx context.Context

	method string
	url    string
	body   any

	Response *Response
	Error    error
}

const ()

func NewRequest(config cdo.Config, httpClient *http.Client, logger *log.Logger, ctx context.Context, method string, url string, body any) *Request {
	return &Request{
		config:     config,
		httpClient: httpClient,
		logger:     logger,

		ctx: ctx,

		method: method,
		url:    url,
		body:   body,
	}
}

// Send wrap send() with retry & delay & timeout... stuff
// TODO: cancel retry when context done
// output: if given, will unmarshal response body into this object, should be a pointer for it to be useful
func (r *Request) Send(output any) error {
	err := retry.Do(func() (bool, error) {

		err := r.send(output)
		if err != nil {
			return false, err
		}
		return true, nil

	}, *retry.NewOptions(
		r.logger,
		r.config.Timeout,
		r.config.Delay,
		r.config.Retries,
		false,
	))

	return err
}

func (r *Request) send(output any) error {
	// clear prev response
	r.Response = nil
	r.Error = nil

	// build net/http.Request
	req, err := r.build()
	if err != nil {
		r.Error = err
		return err
	}

	// send request
	res, err := r.httpClient.Do(req)
	if err != nil {
		r.Error = err
		return err
	}
	defer res.Body.Close()

	// check status
	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		err = fmt.Errorf("failed: code=%d, status=%s, body=%s, readBodyErr=%s, url=%s, method=%s", res.StatusCode, res.Status, string(body), err, r.url, r.method)
		r.Error = err
		return err
	}

	// request is all good, now parse body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		r.Error = err
		return err
	}

	// unmarshal if needed
	if output != nil && len(resBody) > 0 {
		err = json.Unmarshal(resBody, &output)
		if err != nil {
			r.Error = err
			return err
		}
	}

	// set new response
	r.Response = NewResponse(res, resBody, output)
	r.Error = nil

	return nil
}

// build the net/http.Request
func (r *Request) build() (*http.Request, error) {

	bodyReader, err := toReader(r.body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return nil, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}
	r.addAuthHeader(req)
	return req, nil
}

func (r *Request) addAuthHeader(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.config.ApiToken))
	req.Header.Add("Content-Type", "application/json")
}

// toReader try to convert anything to io.Reader.
// Can return nil, which means empty, i.e. empty request body
func toReader(v any) (io.Reader, error) {
	var reader io.Reader
	switch v := v.(type) {
	case io.Reader:
		reader = v
	case string:
		reader = strings.NewReader(v)
	case nil:
		return nil, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}
	return reader, nil
}

func ReadRequestBody[T any](req *http.Request) (*T, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return jsonutil.UnmarshalStruct[T](body)
}
