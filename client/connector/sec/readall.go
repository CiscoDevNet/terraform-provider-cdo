package sec

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type ReadAllInput struct{}

type ReadAllOutput = []ReadOutput

func ReadAll(ctx context.Context, client http.Client, readInp ReadAllInput) (*ReadAllOutput, error) {

	client.Logger.Println("reading all SECs")

	readAllUrl := url.ReadAllSecs(client.BaseUrl())
	readReq := client.NewGet(ctx, readAllUrl)
	var readAllOutput []ReadOutput
	if err := readReq.Send(&readAllOutput); err != nil {
		return nil, err
	}

	return &readAllOutput, nil
}
