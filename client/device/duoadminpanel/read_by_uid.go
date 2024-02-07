package duoadminpanel

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
)

type ReadByUidInput struct {
	Uid string `json:"uid"`
}

func NewReadByUidInput(uid string) ReadByUidInput {
	return ReadByUidInput{
		Uid: uid,
	}
}

type ReadOutput struct {
	Uid   string    `json:"uid"`
	Name  string    `json:"name"`
	State string    `json:"state"`
	Tags  tags.Type `json:"tags"`
}

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	client.Logger.Println("reading duo admin panel")

	readUrl := url.ReadDevice(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var readOutp ReadOutput
	if err := req.Send(&readOutp); err != nil {
		return nil, err
	}

	return &readOutp, nil
}
