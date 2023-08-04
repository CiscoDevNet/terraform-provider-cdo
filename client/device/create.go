package device

import (
	"context"

	"github.com/CiscoDevnet/go-client/internal/http"
	"github.com/CiscoDevnet/go-client/internal/url"
)

type CreateInput struct {
	Name       string `json:"name"`
	DeviceType string `json:"deviceType"`
	LarUid     string `json:"larUid,omitempty"`
	LarType    string `json:"larType"`
	Ipv4       string `json:"ipv4"`
	Model      bool   `json:"model"`

	IgnoreCertifcate bool `json:"ignoreCertificate"`
}

type CreateOutput = ReadOutput

func NewCreateRequestInput(name, deviceType, larUid, larType, ipv4 string, model bool, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:             name,
		DeviceType:       deviceType,
		LarUid:           larUid,
		LarType:          larType,
		Ipv4:             ipv4,
		Model:            model,
		IgnoreCertifcate: ignoreCertificate,
	}
}

func NewCreateRequest(ctx context.Context, client http.Client, createIn CreateInput) *http.Request {

	url := url.CreateDevice(client.BaseUrl())

	req := client.NewPost(ctx, url, createIn)

	return req
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating device")

	req := NewCreateRequest(ctx, client, createInp)

	var outp CreateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
