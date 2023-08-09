package testing

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
)

func MustResponseWithAtMostResponders(method, url string, responders []httpmock.Responder) {
	count := 0

	httpmock.RegisterResponder(
		method,
		url,
		func(r *http.Request) (*http.Response, error) {
			responder := responders[count]
			count += 1
			if count > len(responders) {
				panic(fmt.Sprintf("Too many calls made, method=%s, url=%s, num_expected_calls=%v\n", method, url, len(responders)))
			}

			return responder(r)
		},
	)
}
