package device

import (
	"context"
	"fmt"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type ReadInput struct {
	Uid string `json:"uid"`
}

type ReadByNameInput struct {
	Name string `json:"name"`
}

type ReadByNameAndDeviceTypeInput struct {
	Name string `json:"name"`
	DeviceType string `json:"deviceType"`
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

	IgnoreCertifcate  bool   `json:"ignoreCertificate"`
	ConnectivityState int    `json:"connectivityState,omitempty"`
	ConnectivityError string `json:"connectivityError,omitempty"`
	State             string `json:"state"`
	Status            string `json:"status"`
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

func NewReadByNameAndDeviceTypeRequest(ctx context.Context, client http.Client, readInp ReadByNameAndDeviceTypeInput) *http.Request {

	url := url.ReadDeviceByNameAndDeviceType(client.BaseUrl(), readInp.Name, readInp.DeviceType)

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


func ReadByNameAndDeviceType(ctx context.Context, client http.Client, readInp ReadByNameAndDeviceTypeInput) (*ReadOutput, error) {

	client.Logger.Println("reading Device by name and device type")

	req := NewReadByNameAndDeviceTypeRequest(ctx, client, readInp)

	var arrayOutp []ReadOutput
	if err := req.Send(&arrayOutp); err != nil {
		return nil, err
	}

	if len(arrayOutp) == 0 {
		return nil, fmt.Errorf("no Device by name %s and device type %s found", readInp.Name, readInp.DeviceType)
	}

	if len(arrayOutp) > 1 {
		return nil, fmt.Errorf("multiple devices found with the name: %s and device type: %s", readInp.Name, readInp.DeviceType)
	}

	outp := arrayOutp[0]
	return &outp, nil
}