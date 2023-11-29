package cloudftd

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

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*ReadOutput, error) {

	readUrl := url.ReadDeviceByNameAndType(client.BaseUrl(), readInp.Name, devicetype.CloudFtd)
	req := client.NewGet(ctx, readUrl)

	var readOutp []ReadOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	if len(readOutp) == 0 {
		return nil, fmt.Errorf("%w: cloudftd with name \"%s\" not found", http.NotFoundError, readInp.Name)
	}

	if len(readOutp) > 1 {
		return nil, fmt.Errorf("multiple ftds with name: \"%s\" found, this is unexpected, please report this error at: %s", readInp.Name, cdo.TerraformProviderCDOIssuesUrl)
	}

	return &readOutp[0], nil
}
