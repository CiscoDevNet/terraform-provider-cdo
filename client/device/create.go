package device

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type CreateInput struct {
	// required
	Name       string          `json:"name"`
	DeviceType devicetype.Type `json:"deviceType"`
	Model      bool            `json:"model"`

	// optional
	ConnectorUid       string       `json:"larUid,omitempty"`
	ConnectorType      string       `json:"larType,omitempty"`
	SocketAddress      string       `json:"ipv4,omitempty"`
	IgnoreCertificate  *bool        `json:"ignoreCertificate,omitempty"`
	Metadata           *interface{} `json:"metadata,omitempty"`
	Tags               tags.Type    `json:"tags,omitempty"`
	EnableOobDetection *bool        `json:"enableOobDetection,omitempty"`
}

type CreateOutput = ReadOutput

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
