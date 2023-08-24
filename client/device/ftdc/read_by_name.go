package ftdc

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
)

type ReadByNameInput struct {
	Name string
}

func NewReadByNameInput(name string) ReadByNameInput {
	return ReadByNameInput{
		Name: name,
	}
}

type ReadByNameOutput struct {
	Uid      string   `json:"uid"`
	Name     string   `json:"name"`
	Metadata Metadata `json:"metadata,omitempty"`
}

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*ReadByNameOutput, error) {

	readUrl := url.ReadDeviceByNameAndType(client.BaseUrl(), readInp.Name, devicetype.Ftdc)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadByNameOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
