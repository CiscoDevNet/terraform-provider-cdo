package device

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type ReadInput struct {
	Uid string `json:"uid"`
}

type ReadOutput struct {
	Uid             string `json:"uid"`
	Name            string `json:"name"`
	CreatedDate     int64  `json:"createdDate"`
	LastUpdatedDate int64  `json:"lastUpdatedDate"`
	DeviceType      string `json:"deviceType"`
	LarUid          string `json:"larUid"`
	LarType         string `json:"larType"`
	Ipv4            string `json:"ipv4"`
	Port            string `json:"port"`
	Host            string `json:"host"`
}

func NewReadInput(uid string) *ReadInput {
	return &ReadInput{
		Uid: uid,
	}
}

func NewReadRequest(ctx context.Context, client http.Client, readInp ReadInput) *http.Request {

	url := url.ReadDevice(client.BaseUrl(), readInp.Uid)

	req := client.NewGet(ctx, url)

	return req
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading device")

	req := NewReadRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
