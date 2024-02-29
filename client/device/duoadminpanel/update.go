package duoadminpanel

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
)

type UpdateInput struct {
	Uid  string    `json:"-"`
	Name string    `json:"name"`
	Tags tags.Type `json:"tags"`
}

type UpdateOutput = ReadOutput

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating duo admin panel")

	updateUrl := url.UpdateDuoAdminPanel(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, updateUrl, updateInp)

	var updateOutput UpdateOutput
	if err := req.Send(&updateOutput); err != nil {
		return nil, err
	}

	return &updateOutput, nil

}
