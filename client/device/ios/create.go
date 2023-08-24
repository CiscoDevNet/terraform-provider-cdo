package ios

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
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
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

func NewCreateRequestInput(name, connectorUid, connectorType, socketAddress, username, password string, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:              name,
		ConnectorUid:      connectorUid,
		ConnectorType:     connectorType,
		SocketAddress:     socketAddress,
		Username:          username,
		Password:          password,
		IgnoreCertificate: ignoreCertificate,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating ios device")

	deviceCreateOutp, err := device.Create(ctx, client, *device.NewCreateRequestInput(
		createInp.Name, "IOS", createInp.ConnectorUid, createInp.ConnectorType, createInp.SocketAddress, false, createInp.IgnoreCertificate,
	))
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

	err = retry.Do(iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, state.PRE_READ_METADATA), *retry.NewOptionsWithLogger(client.Logger))
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

	err = retry.Do(iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, state.DONE), *retry.NewOptionsWithLogger(client.Logger))
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
	}
	return &createOutp, nil
}
