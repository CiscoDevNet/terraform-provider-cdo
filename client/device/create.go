package device

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type CreateInput struct {
	Name          string `json:"name"`
	DeviceType    string `json:"deviceType"`
	ConnectorUid  string `json:"larUid,omitempty"`
	ConnectorType string `json:"larType"`
	SocketAddress string `json:"ipv4"`
	Model         bool   `json:"model"`

	IgnoreCertificate bool `json:"ignoreCertificate"`
}

type CreateOutput = ReadOutput

func NewCreateRequestInput(name, deviceType, connectorUid, connectorType, socketAddress string, model bool, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:              name,
		DeviceType:        deviceType,
		ConnectorUid:      connectorUid,
		ConnectorType:     connectorType,
		SocketAddress:     socketAddress,
		Model:             model,
		IgnoreCertificate: ignoreCertificate,
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
