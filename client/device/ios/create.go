package ios

import (
	"context"
	"fmt"
	"strings"

	"github.com/CiscoDevnet/go-client/connector/sdc"
	"github.com/CiscoDevnet/go-client/device"
	"github.com/CiscoDevnet/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/go-client/internal/http"
	"github.com/CiscoDevnet/go-client/internal/retry"
)

type CreateInput struct {
	Name    string
	LarUid  string
	LarType string
	Ipv4    string

	Username string
	Password string

	IgnoreCertifcate bool
}

type CreateOutput struct {
	Uid        string `json:"uid"`
	Name       string `json:"Name"`
	DeviceType string `json:"deviceType"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Ipv4       string `json:"ipv4"`
	LarType    string `json:"larType"`
	LarUid     string `json:"larUid"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

const (
	IosStatePreReadMetadata = "$PRE_READ_METADATA"
	IosStateDone            = "DONE"
)

func NewCreateRequestInput(name, larUid, larType, ipv4, username, password string, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:             name,
		LarUid:           larUid,
		LarType:          larType,
		Ipv4:             ipv4,
		Username:         username,
		Password:         password,
		IgnoreCertifcate: ignoreCertificate,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating ios device")

	deviceCreateOutp, err := device.Create(ctx, client, *device.NewCreateRequestInput(
		createInp.Name, "IOS", createInp.LarUid, createInp.LarType, createInp.Ipv4, false, createInp.IgnoreCertifcate,
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
		larReadRes, err := sdc.ReadByUid(ctx, client, sdc.ReadInput{
			LarUid: deviceCreateOutp.LarUid,
		})
		if err != nil {
			return nil, &CreateError{
				Err:               err,
				CreatedResourceId: createdResourceId,
			}
		}
		publicKey = &larReadRes.PublicKey
	}

	err = retry.Do(iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, IosStatePreReadMetadata), *retry.NewOptionsWithLogger(client.Logger))
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

	err = retry.Do(iosconfig.UntilState(ctx, client, deviceCreateOutp.Uid, IosStateDone), *retry.NewOptionsWithLogger(client.Logger))
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
		LarUid:     deviceCreateOutp.LarUid,
		LarType:    deviceCreateOutp.LarType,
	}
	return &createOutp, nil
}
