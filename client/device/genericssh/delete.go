package genericssh

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type DeleteInput struct {
	Uid string
}

type DeleteOutput = device.ReadOutput

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	client.Logger.Println("deleting generic ssh")

	_, err := device.Delete(ctx, client, *device.NewDeleteInput(deleteInp.Uid))
	if err != nil {
		return nil, err
	}

	return &DeleteOutput{}, nil
}
