package asa

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type ReadInput = device.ReadInput

type ReadOutput struct {
	Uid             string `json:"uid"`
	Name            string `json:"name"`
	CreatedDate     int64  `json:"createdDate"`
	LastUpdatedDate int64  `json:"lastUpdatedDate"`
	DeviceType      string `json:"deviceType"`
	ConnectorUid    string `json:"larUid"`
	ConnectorType   string `json:"larType"`
	SocketAddress   string `json:"ipv4"`
	Port            string `json:"port"`
	Host            string `json:"host"`

	IgnoreCertifcate  bool   `json:"ignoreCertificate"`
	ConnectivityState int    `json:"connectivityState,omitempty"`
	ConnectivityError string `json:"connectivityError,omitempty"`
	State             string `json:"state"`
	Status            string `json:"status"`
}

func NewReadInput(uid string) *ReadInput {
	return device.NewReadInput(uid)
}

func NewReadRequest(ctx context.Context, client http.Client, readInp ReadInput) *http.Request {
	return device.NewReadRequest(ctx, client, readInp)
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading asa device")

	req := NewReadRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
