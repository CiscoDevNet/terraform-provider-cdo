package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadInput struct {
	Uid string
}

type ReadOutput struct {
	Name            string `json:"name"`
	Uid             string `json:"uid"`
	BootStrapData   string `json:"bootstrapData"`
	TokenExpiryTime int64  `json:"tokenExpiryTime"`
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	readUrl := url.ReadSec(client.BaseUrl(), readInp.Uid)
	readReq := client.NewGet(ctx, readUrl)
	var readOutput ReadOutput
	if err := readReq.Send(&readOutput); err != nil {
		return nil, err
	}

	return &readOutput, nil
}
