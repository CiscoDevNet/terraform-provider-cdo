package asa

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

type Metadata struct {
	IsNewPolicyObjectModel string `json:"isNewPolicyObjectModel"` // yes it is a string, but it should be either "true" or "false" :/
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

type createBody struct {
	Name              string               `json:"name"`
	DeviceAddress     string               `json:"deviceAddress"`
	Username          string               `json:"username"`
	Password          string               `json:"password"`
	ConnectorType     string               `json:"connectorType"`
	IgnoreCertificate bool                 `json:"ignoreCertificate"`
	ConnectorName     string               `json:"connectorName"`
	Labels            publicapilabels.Type `json:"labels"`
}

func (r *CreateError) Error() string {
	return r.Err.Error()
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

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating asa")

	createUrl := url.CreateAsa(client.BaseUrl())

	var connectorName string
	if createInp.ConnectorType == "SDC" {
		conn, err := connector.ReadByUid(ctx, client, connector.ReadByUidInput{ConnectorUid: createInp.ConnectorUid})
		if err != nil {
			return nil, &CreateError{
				Err:               err,
				CreatedResourceId: nil,
			}
		}
		connectorName = conn.Name
	}

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createBody{
			Name:              createInp.Name,
			DeviceAddress:     createInp.SocketAddress,
			Username:          createInp.Username,
			Password:          createInp.Password,
			ConnectorType:     createInp.ConnectorType,
			IgnoreCertificate: createInp.IgnoreCertificate,
			ConnectorName:     connectorName,
			Labels:            createInp.Labels,
		},
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}
	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for Asa to onboard...",
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

	readOut, err := Read(ctx, client, ReadInput{Uid: transaction.EntityUid})
	if err == nil {
		return readOut, nil
	} else {
		return readOut, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}
}
