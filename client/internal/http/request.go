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
	netUrl "net/url"
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

	Header      http.Header
	QueryParams netUrl.Values

	Response *Response
	Error    error
}

func NewRequest(config cdo.Config, httpClient *http.Client, logger *log.Logger, ctx context.Context, method string, url string, body any) *Request {
	return &Request{
		config:     config,
		httpClient: httpClient,
		logger:     logger,

		ctx: ctx,

		method: method,
		url:    url,
		body:   body,

		Header:      make(http.Header),
		QueryParams: make(netUrl.Values),
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
		err = fmt.Errorf("failed: url=%s, code=%d, status=%s, body=%s, readBodyErr=%s, method=%s, header=%s", r.url, res.StatusCode, res.Status, string(body), err, r.method, r.Header)
		r.Error = err
		return err
	}

	// request is all good, now parse body
	resBody, err := io.ReadAll(res.Body)
	fmt.Printf("\n\nsuccess: url=%s, code=%d, status=%s, body=%s, readBodyErr=%s, method=%s, header=%s, queryParams=%s\n", r.url, res.StatusCode, res.Status, string(resBody), err, r.method, r.Header, r.QueryParams)
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

	// TODO: remove these debug lines
	//if r.method != "GET" && r.method != "DELETE" {
	//	bodyReader2, err := toReader(r.body)
	//	if err != nil {
	//		return nil, err
	//	}
	//	bs, err := io.ReadAll(bodyReader2)
	//	if err != nil {
	//		return nil, err
	//	}
	//	fmt.Println("request_check")
	//	fmt.Printf("Request: %+v\n", r)
	//	fmt.Printf("Request: %s, %s, %s\n", r.url, r.method, string(bs))
	//}

	req, err := http.NewRequest(r.method, r.url, bodyReader)
	if err != nil {
		return nil, err
	}
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	r.addHeaders(req)
	r.addQueryParams(req)
	return req, nil
}

func (r *Request) addQueryParams(req *http.Request) {
	q := req.URL.Query()
	for k, vs := range r.QueryParams {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	s := q.Encode()
	if s != "" {
		fmt.Printf("\n\nencoded_query=%s\n\n", s)
	}
	req.URL.RawQuery = s
}

func (r *Request) addHeaders(req *http.Request) {
	r.addAuthHeader(req)
	r.addOtherHeader(req)
}

func (r *Request) addAuthHeader(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.config.ApiToken))
	req.Header.Add("Content-Type", "application/json")
}

func (r *Request) addOtherHeader(req *http.Request) {
	for k, vs := range r.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
}

// toReader tries to convert anything to io.Reader.
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
