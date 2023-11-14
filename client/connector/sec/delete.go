package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type DeleteInput struct {
	Uid string
}

type DeleteOutput struct {
}

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	deleteUrl := url.DeleteSec(client.BaseUrl(), deleteInp.Uid)
	deleteReq := client.NewDelete(ctx, deleteUrl)
	var deleteOutput DeleteOutput
	if err := deleteReq.Send(&deleteOutput); err != nil {
		return nil, err
	}

	return &deleteOutput, nil
}
