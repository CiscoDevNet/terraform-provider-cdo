package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type UpdateInput struct {
	Uid  string `json:"-"`
	Name string `json:"name"`
}

func NewUpdateInput(uid string, name string) UpdateInput {
	return UpdateInput{
		Uid:  uid,
		Name: name,
	}
}

type UpdateSdcOutput = UpdateOutput

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateSdcOutput, error) {

	url := url.UpdateSdc(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	var updateOutp UpdateSdcOutput
	if err := req.Send(&updateOutp); err != nil {
		return &UpdateSdcOutput{}, nil
	}

	return &updateOutp, nil
}
