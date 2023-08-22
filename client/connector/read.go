package connector

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadByUidInput struct {
	ConnectorUid string
}

type ReadByNameInput struct {
	ConnectorName string
}

type ReadOutput struct {
	Uid              string          `json:"uid"`
	Name             string          `json:"name"`
	DefaultConnector bool            `json:"defaultLar"`
	Cdg              bool            `json:"cdg"`
	TenantUid        string          `json:"tenantUid"`
	PublicKey        model.PublicKey `json:"larPublicKey"`
}

func NewReadByUidInput(connectorUid string) *ReadByUidInput {
	return &ReadByUidInput{
		ConnectorUid: connectorUid,
	}
}

func NewReadByNameInput(ConnectorName string) *ReadByNameInput {
	return &ReadByNameInput{
		ConnectorName: ConnectorName,
	}
}

func newReadByUidRequest(ctx context.Context, client http.Client, readInp ReadByUidInput) *http.Request {

	url := url.ReadConnectorByUid(client.BaseUrl(), readInp.ConnectorUid)

	req := client.NewGet(ctx, url)

	return req
}

func newReadByNameRequest(ctx context.Context, client http.Client, readInp ReadByNameInput) *http.Request {

	url := url.ReadConnectorByName(client.BaseUrl(), readInp.ConnectorName)

	req := client.NewGet(ctx, url)

	return req
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	client.Logger.Println("reading connector by uid")

	req := newReadByUidRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*ReadOutput, error) {

	client.Logger.Println("reading connector by name")

	req := newReadByNameRequest(ctx, client, readInp)

	var arrayOutp []ReadOutput
	if err := req.Send(&arrayOutp); err != nil {
		return nil, err
	}

	if len(arrayOutp) == 0 {
		return nil, fmt.Errorf("no connector found")
	}

	if len(arrayOutp) > 1 {
		return nil, fmt.Errorf("multiple connectors found with the name: %s", readInp.ConnectorName)
	}

	outp := arrayOutp[0]
	return &outp, nil
}
