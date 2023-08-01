package sdc

import (
	"context"
	"fmt"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type ReadInput struct {
	LarUid string
}

type ReadByNameInput struct {
	LarName string
}

type ReadOutput struct {
	Uid        string    `json:"uid"`
	Name       string    `json:"name"`
	DefaultLar bool      `json:"defaultLar"`
	Cdg        bool      `json:"cdg"`
	TenantUid  string    `json:"tenantUid"`
	PublicKey  PublicKey `json:"larPublicKey"`
}

type PublicKey struct {
	EncodedKey string `json:"encodedKey"`
	Version    int64  `json:"version"`
	KeyId      string `json:"keyId"`
}

func NewReadInput(larUid string) *ReadInput {
	return &ReadInput{
		LarUid: larUid,
	}
}

func NewReadByNameInput(larName string) *ReadByNameInput {
	return &ReadByNameInput{
		LarName: larName,
	}
}

func newReadByUidRequest(ctx context.Context, client http.Client, readInp ReadInput) *http.Request {

	url := url.ReadSdcByUid(client.BaseUrl(), readInp.LarUid)

	req := client.NewGet(ctx, url)

	return req
}

func newReadByNameRequest(ctx context.Context, client http.Client, readInp ReadByNameInput) *http.Request {

	url := url.ReadSdcByName(client.BaseUrl(), readInp.LarName)

	req := client.NewGet(ctx, url)

	return req
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading sdc")

	req := newReadByUidRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*ReadOutput, error) {

	client.Logger.Println("reading sdc by name")

	req := newReadByNameRequest(ctx, client, readInp)

	var arrayOutp []ReadOutput
	if err := req.Send(&arrayOutp); err != nil {
		return nil, err
	}

	if len(arrayOutp) == 0 {
		return nil, fmt.Errorf("no SDC found")
	}

	if len(arrayOutp) > 1 {
		return nil, fmt.Errorf("multiple SDCs found with the name: %s", readInp.LarName)
	}

	outp := arrayOutp[0]
	return &outp, nil
}
