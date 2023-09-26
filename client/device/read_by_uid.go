package device

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

type ReadOutput struct {
	Uid             string          `json:"uid"`
	Name            string          `json:"name"`
	CreatedDate     int64           `json:"createdDate"`
	LastUpdatedDate int64           `json:"lastUpdatedDate"`
	DeviceType      devicetype.Type `json:"deviceType"`
	ConnectorUid    string          `json:"larUid"`
	ConnectorType   string          `json:"larType"`
	SocketAddress   string          `json:"ipv4"`
	SoftwareVersion string          `json:"softwareVersion"`
	Port            string          `json:"port"`
	Host            string          `json:"host"`
	Tags            tags.Type       `json:"tags"`

	IgnoreCertificate bool   `json:"ignoreCertificate"`
	ConnectivityState int    `json:"connectivityState,omitempty"`
	ConnectivityError string `json:"connectivityError,omitempty"`
	State             string `json:"state"`
	Status            string `json:"status"`
}

func NewReadByUidInput(uid string) *ReadByUidInput {
	return &ReadByUidInput{
		Uid: uid,
	}
}

func NewReadByUidRequest(ctx context.Context, client http.Client, readInp ReadByUidInput) *http.Request {
	readUrl := url.ReadDevice(client.BaseUrl(), readInp.Uid)

	req := client.NewGet(ctx, readUrl)

	return req
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	client.Logger.Println("reading device")

	req := NewReadByUidRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
