package ios

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/publicapilabels"
)

type CreateInput struct {
	Name          string
	ConnectorUid  string
	ConnectorType string
	SocketAddress string
	Labels        publicapilabels.Type

	Username string
	Password string

	IgnoreCertificate bool
}

type CreateOutput = ReadOutput

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

type createBody struct {
	Name          string `json:"name"`
	ConnectorName string `json:"connectorName"`
	ConnectorType string `json:"connectorType"`
	SocketAddress string `json:"deviceAddress"`

	Username string `json:"username"`
	Password string `json:"password"`

	IgnoreCertificate bool                 `json:"ignoreCertificate"`
	Labels            publicapilabels.Type `json:"labels"`
}

func NewCreateRequestInput(name, connectorUid, connectorType, socketAddress, username, password string, ignoreCertificate bool, labels publicapilabels.Type) *CreateInput {
	return &CreateInput{
		Name:              name,
		ConnectorUid:      connectorUid,
		ConnectorType:     connectorType,
		SocketAddress:     socketAddress,
		Username:          username,
		Password:          password,
		IgnoreCertificate: ignoreCertificate,
		Labels:            labels,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating ios device")

	createUrl := url.CreateIos(client.BaseUrl())

	conn, err := connector.ReadByUid(ctx, client, *connector.NewReadByUidInput(createInp.ConnectorUid))
	if err != nil {
		return nil, err
	}

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createBody{
			Name:              createInp.Name,
			ConnectorName:     conn.Name,
			ConnectorType:     createInp.ConnectorType,
			SocketAddress:     createInp.SocketAddress,
			Username:          createInp.Username,
			Password:          createInp.Password,
			IgnoreCertificate: createInp.IgnoreCertificate,
			Labels:            createInp.Labels,
		},
	)
	if err != nil {
		_, _ = Delete(ctx, client, *NewDeleteInput(transaction.EntityUid))
		return nil, err
	}
	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for IOS to onboard...",
	)
	if err != nil {
		_, _ = Delete(ctx, client, *NewDeleteInput(transaction.EntityUid))
		return nil, err
	}

	return Read(ctx, client, *NewReadInput(transaction.EntityUid))
}
