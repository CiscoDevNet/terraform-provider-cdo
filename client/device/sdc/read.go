package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type ReadInput struct {
	LarUid string
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

func NewReadRequest(ctx context.Context, client http.Client, readInp ReadInput) *http.Request {

	url := url.ReadSdc(client.BaseUrl(), readInp.LarUid)

	req := client.NewGet(ctx, url)

	return req
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading sdc")

	req := NewReadRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
