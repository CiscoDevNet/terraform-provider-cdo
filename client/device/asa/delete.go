package asa

import (
	"context"

	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/internal/http"
)

type DeleteInput = device.DeleteInput
type DeleteOutput = device.DeleteOutput

func NewDeleteInput(uid string) *DeleteInput {
	return &DeleteInput{
		Uid: uid,
	}
}

func NewDeleteRequest(ctx context.Context, client http.Client, deleteInp DeleteInput) *http.Request {
	return device.NewDeleteRequest(ctx, client, deleteInp)
}

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	client.Logger.Println("deleting asa device")

	return device.Delete(ctx, client, deleteInp)
}
