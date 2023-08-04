package sdc

import (
	"context"

	"github.com/CiscoDevnet/go-client/internal/http"
	"github.com/CiscoDevnet/go-client/internal/url"
)

type ReadAllInput struct{}

type ReadAllOutput = []ReadOutput

func NewReadAllInput() *ReadAllInput {
	return &ReadAllInput{}
}

func NewReadAllRequest(ctx context.Context, client http.Client, readAllInp ReadAllInput) *http.Request {

	url := url.ReadAllSdcs(client.BaseUrl())

	req := client.NewGet(ctx, url)

	return req
}

// TODO: Change the return type to return value type over pointer (*ReadAllOutput -> ReadAllOutput). Slices are references in golang.
func ReadAll(ctx context.Context, client http.Client, readAllInp ReadAllInput) (*ReadAllOutput, error) {

	client.Logger.Println("reading all sdcs")

	req := NewReadAllRequest(ctx, client, readAllInp)

	var outp ReadAllOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
