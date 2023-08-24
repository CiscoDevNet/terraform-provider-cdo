package ftdc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type DeleteInput struct {
	Uid string
}

func NewDeleteInput(uid string) DeleteInput {
	return DeleteInput{
		Uid: uid,
	}
}

type DeleteOutput = ReadByUidOutput

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {
	
	// TODO: is this all it takes to delete a ftdc? What about the underlying virtual ftd?
	deleteUrl := url.DeleteDevice(client.BaseUrl(), deleteInp.Uid)
	req := client.NewDelete(ctx, deleteUrl)
	var deleteOutp DeleteOutput
	if err := req.Send(&deleteOutp); err != nil {
		return nil, err
	}

	return &deleteOutp, nil

}
