package ftdc

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
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

	var readOutp []ReadByNameOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	if len(readOutp) == 0 {
		return nil, fmt.Errorf("ftd with name: \"%s\" not found", readInp.Name)
	}

	if len(readOutp) > 1 {
		return nil, fmt.Errorf("multiple ftds with name: \"%s\" found, this is unexpected, please report this error at: %s", readInp.Name, cdo.TerraformProviderCDOIssuesUrl)
	}

	return &readOutp[0], nil
}
