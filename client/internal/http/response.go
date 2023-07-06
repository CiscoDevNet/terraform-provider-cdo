package http

import "net/http"

// Represent response of Request.
// Raw: raw response
// Body: response body bytes
// Output: unmarshalled object, if any
type Response struct {
	Raw    *http.Response
	Body   []byte
	Output any
}

func NewResponse(res *http.Response, resBody []byte, resObj any) *Response {
	return &Response{
		Raw:    res,
		Body:   resBody,
		Output: resObj,
	}
}
