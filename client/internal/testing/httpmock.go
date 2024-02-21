package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func GetCallCount(method, path string) (int, bool) {
	callCounts := httpmock.GetCallCountInfo()

	callCount, ok := callCounts[fmt.Sprintf("%s %s", method, path)]
	return callCount, ok
}

func AssertEndpointCalledTimes(method, path string, times int, t *testing.T) {
	actualCount, found := GetCallCount(method, path)

	if !found {
		t.Errorf("%s %s was not called!", method, path)
		return
	}

	if actualCount != times {
		t.Errorf("expected %s %s to be called %d times, but was called: %d times", method, path, times, actualCount)
	}
}

func MockPostAccepted(url string, body any) {
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewJsonResponderOrPanic(202, body))
}

func MockPostError(url string, body any) {
	httpmock.RegisterResponder(http.MethodPost, url, httpmock.NewJsonResponderOrPanic(500, body))
}

func MockGetOk(url string, body any) {
	httpmock.RegisterResponder(http.MethodGet, url, httpmock.NewJsonResponderOrPanic(200, body))
}

func MockGetError(url string, body any) {
	httpmock.RegisterResponder(http.MethodGet, url, httpmock.NewJsonResponderOrPanic(500, body))
}
