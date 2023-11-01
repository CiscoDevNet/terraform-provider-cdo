package ios

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/devicetype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

type CreateInput struct {
	Name          string
	ConnectorUid  string
	ConnectorType string
	SocketAddress string
	Tags          tags.Type

	Username string
	Password string

	IgnoreCertificate bool
}

type CreateOutput struct {
	Uid           string          `json:"uid"`
	Name          string          `json:"Name"`
	DeviceType    devicetype.Type `json:"deviceType"`
	Host          string          `json:"host"`
	Port          string          `json:"port"`
	SocketAddress string          `json:"ipv4"`
	ConnectorType string          `json:"larType"`
	ConnectorUid  string          `json:"larUid"`
	Tags          tags.Type       `json:"tags"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

func NewCreateRequestInput(name, connectorUid, connectorType, socketAddress, username, password string, ignoreCertificate bool, tags tags.Type) *CreateInput {
	return &CreateInput{
		Name:              name,
		ConnectorUid:      connectorUid,
		ConnectorType:     connectorType,
		SocketAddress:     socketAddress,
		Username:          username,
		Password:          password,
		IgnoreCertificate: ignoreCertificate,
		Tags:              tags,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating ios device")

	deviceCreateOutp, err := device.Create(
		ctx,
		client,
		device.NewCreateInputBuilder().
			Name(createInp.Name).
			DeviceType(devicetype.Ios).
			ConnectorUid(createInp.ConnectorUid).
			ConnectorType(createInp.ConnectorType).
			SocketAddress(createInp.SocketAddress).
			Model(false).
			IgnoreCertificate(&createInp.IgnoreCertificate).
			Metadata(nil).
			Tags(createInp.Tags).
			Build(),
	)

	var createdResourceId *string = nil
	if deviceCreateOutp != nil {
		createdResourceId = &deviceCreateOutp.Uid
	}

	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: createdResourceId,
		}
	}

	// encrypt credentials on prem connector
	var publicKey *model.PublicKey
	if strings.EqualFold(deviceCreateOutp.ConnectorType, "SDC") {

		// on-prem connector requires encryption
		client.Logger.Println("decrypting public key from connector for encryption")

		if deviceCreateOutp.ConnectorUid == "" {
			return nil, &CreateError{
				Err:               fmt.Errorf("connector uid not found"),
				CreatedResourceId: createdResourceId,
			}

		}

		// read connector public key
		connectorReadRes, err := connector.ReadByUid(ctx, client, *connector.NewReadByUidInput(deviceCreateOutp.ConnectorUid))
		if err != nil {
			return nil, &CreateError{
				Err:               err,
				CreatedResourceId: createdResourceId,
			}
		}
		publicKey = &connectorReadRes.PublicKey
	}

	err = retry.Do(
		ctx,
		iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, state.PRE_READ_METADATA),
		retry.NewOptionsBuilder().
			Message("Waiting for IOS device to be onboarded to CDO...").
			Retries(retry.DefaultRetries).
			Delay(retry.DefaultDelay).
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(retry.DefaultTimeout).
			Build(),
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: createdResourceId,
		}
	}

	// update ios config credentials
	client.Logger.Println("updating ios config credentials")

	iosConfigUpdateInp := iosconfig.NewUpdateInput(
		deviceCreateOutp.Uid,
		createInp.Username,
		createInp.Password,
		publicKey,
	)
	_, err = iosconfig.Update(ctx, client, *iosConfigUpdateInp)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: createdResourceId,
		}
	}

	// poll until ios config state done
	client.Logger.Println("waiting for device to reach state done")

	err = retry.Do(
		ctx,
		iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, state.DONE),
		retry.NewOptionsBuilder().
			Message("Waiting for IOS device to be onboarded to CDO...").
			Retries(retry.DefaultRetries).
			Delay(retry.DefaultDelay).
			Logger(client.Logger).
			EarlyExitOnError(true).
			Timeout(retry.DefaultTimeout).
			Build(),
	)
	if err != nil {
		return nil, &CreateError{
			Err:               err,
			CreatedResourceId: createdResourceId,
		}
	}

	// done!
	createOutp := CreateOutput{
		Uid:           deviceCreateOutp.Uid,
		Name:          deviceCreateOutp.Name,
		DeviceType:    deviceCreateOutp.DeviceType,
		Host:          deviceCreateOutp.Host,
		Port:          deviceCreateOutp.Port,
		SocketAddress: deviceCreateOutp.SocketAddress,
		ConnectorUid:  deviceCreateOutp.ConnectorUid,
		ConnectorType: deviceCreateOutp.ConnectorType,
		Tags:          deviceCreateOutp.Tags,
	}
	return &createOutp, nil
}
