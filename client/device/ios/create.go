package ios

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

type CreateInput struct {
	Name    string
	SdcUid  string
	SdcType string
	Ipv4    string

	Username string
	Password string

	IgnoreCertificate bool
}

type CreateOutput struct {
	Uid        string `json:"uid"`
	Name       string `json:"Name"`
	DeviceType string `json:"deviceType"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Ipv4       string `json:"ipv4"`
	SdcType    string `json:"larType"`
	SdcUid     string `json:"larUid"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

func NewCreateRequestInput(name, sdcUid, sdcType, ipv4, username, password string, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:              name,
		SdcUid:            sdcUid,
		SdcType:           sdcType,
		Ipv4:              ipv4,
		Username:          username,
		Password:          password,
		IgnoreCertificate: ignoreCertificate,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating ios device")

	deviceCreateOutp, err := device.Create(ctx, client, *device.NewCreateRequestInput(
		createInp.Name, "IOS", createInp.SdcUid, createInp.SdcType, createInp.Ipv4, false, createInp.IgnoreCertificate,
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

	// encrypt credentials for SDC on prem lar
	var publicKey *sdc.PublicKey
	if strings.EqualFold(deviceCreateOutp.LarType, "SDC") {

		// on-prem lar requires encryption
		client.Logger.Println("decrypting public key from sdc for encrpytion")

		if deviceCreateOutp.LarUid == "" {
			return nil, &CreateError{
				Err:               fmt.Errorf("sdc uid not found"),
				CreatedResourceId: createdResourceId,
			}

		}

		// read lar public key
		larReadRes, err := sdc.ReadByUid(ctx, client, sdc.ReadByUidInput{
			SdcUid: deviceCreateOutp.LarUid,
		})
		if err != nil {
			return nil, &CreateError{
				Err:               err,
				CreatedResourceId: createdResourceId,
			}
		}
		publicKey = &larReadRes.PublicKey
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
		Uid:        deviceCreateOutp.Uid,
		Name:       deviceCreateOutp.Name,
		DeviceType: deviceCreateOutp.DeviceType,
		Host:       deviceCreateOutp.Host,
		Port:       deviceCreateOutp.Port,
		Ipv4:       deviceCreateOutp.Ipv4,
		SdcUid:     deviceCreateOutp.LarUid,
		SdcType:    deviceCreateOutp.LarType,
	}
	return &createOutp, nil
}
