package asa

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	SoftwareVersion   string
	AsdmVersion       string
}

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

func NewCreateRequestInput(name, connectorUid, connectorType, socketAddress, username, password string, ignoreCertificate bool, labels publicapilabels.Type, softwareVersion string, asdmVersion string) *CreateInput {
	return &CreateInput{
		Name:              name,
		ConnectorUid:      connectorUid,
		ConnectorType:     connectorType,
		SocketAddress:     socketAddress,
		Username:          username,
		Password:          password,
		IgnoreCertificate: ignoreCertificate,
		Labels:            labels,
		SoftwareVersion:   softwareVersion,
		AsdmVersion:       asdmVersion,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*ReadOutput, *ReadSpecificOutput, *CreateError) {

	client.Logger.Println("creating asa device")

	createUrl := url.CreateAsa(client.BaseUrl())

	var connectorName string
	if createInp.ConnectorType == "SDC" {
		conn, err := connector.ReadByUid(ctx, client, connector.ReadByUidInput{ConnectorUid: createInp.ConnectorUid})
		if err != nil {
			return nil, nil, &CreateError{
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
		return nil, nil, &CreateError{
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
		return nil, nil, &CreateError{
			Err:               err,
			CreatedResourceId: &transaction.EntityUid,
		}
	}

	readOut, err := Read(ctx, client, ReadInput{Uid: transaction.EntityUid})
	if err == nil {
		if err = validateSoftwareVersion(readOut, createInp.SoftwareVersion); err != nil {
			client.Logger.Println("Version validation failed...")
			return nil, nil, &CreateError{
				Err:               err,
				CreatedResourceId: &transaction.EntityUid,
			}
		}

		readSpecificOut, err := ReadSpecific(ctx, client, ReadSpecificInput{Uid: transaction.EntityUid})
		if err == nil {
			if err = validateAsdmVersion(readSpecificOut, createInp.AsdmVersion); err != nil {
				client.Logger.Println("ASDM version validation failed...")
				return nil, nil, &CreateError{
					Err:               err,
					CreatedResourceId: &transaction.EntityUid,
				}
			}
			return readOut, readSpecificOut, nil
		}
	}

	return nil, nil, &CreateError{
		Err:               err,
		CreatedResourceId: &transaction.EntityUid,
	}
}

func validateSoftwareVersion(readOut *ReadOutput, softwareVersion string) error {
	if strings.TrimSpace(softwareVersion) != "" && readOut.SoftwareVersion != softwareVersion {
		return errors.New(fmt.Sprintf("ASA Software version mismatch. Specified software version %s does not match actual software version %s on ASA device", softwareVersion, readOut.SoftwareVersion))
	}

	return nil
}

func validateAsdmVersion(readOut *ReadSpecificOutput, asdmVersion string) error {
	if strings.TrimSpace(asdmVersion) != "" && readOut.Metadata.AsdmVersion != asdmVersion {
		return errors.New(fmt.Sprintf("ASDM version mismatch. Specified ASDM version %s does not match actual ASDM version %s on ASA device", asdmVersion, readOut.Metadata.AsdmVersion))
	}

	return nil
}
