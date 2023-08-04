package sdc

import (
	"context"

	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type DeleteInput struct {
	Uid string `json:"-"`
}

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

type DeleteOutput struct {
}

func Delete(ctx context.Context, client http.Client, inp DeleteInput) (*DeleteOutput, error) {

	url := url.DeleteSdc(client.BaseUrl(), inp.Uid)

	req := client.NewDelete(ctx, url)

	var deleteOutp DeleteOutput
	if err := req.Send(&deleteOutp); err != nil {
		return &DeleteOutput{}, err
	}

	return &deleteOutp, nil
}
