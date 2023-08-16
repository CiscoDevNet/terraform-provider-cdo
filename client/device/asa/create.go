package asa

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"
)

type CreateInput struct {
	Name          string
	SdcUid        string
	SdcType       string
	SocketAddress string

	Username string
	Password string

	IgnoreCertifcate bool
}

type CreateOutput struct {
	Uid           string `json:"uid"`
	Name          string `json:"Name"`
	DeviceType    string `json:"deviceType"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	SocketAddress string `json:"ipv4"`
	SdcType       string `json:"larType"`
	SdcUid        string `json:"larUid"`
}

type CreateError struct {
	Err               error
	CreatedResourceId *string
}

func (r *CreateError) Error() string {
	return r.Err.Error()
}

func NewCreateRequestInput(name, larUid, larType, ipv4, username, password string, ignoreCertificate bool) *CreateInput {
	return &CreateInput{
		Name:             name,
		SdcUid:           larUid,
		SdcType:          larType,
		SocketAddress:    ipv4,
		Username:         username,
		Password:         password,
		IgnoreCertifcate: ignoreCertificate,
	}
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, *CreateError) {

	client.Logger.Println("creating asa device")

	deviceCreateOutp, err := device.Create(ctx, client, *device.NewCreateRequestInput(
		createInp.Name, "ASA", createInp.SdcUid, createInp.SdcType, createInp.SocketAddress, false, createInp.IgnoreCertifcate,
	))
	var createdResourceId *string = nil
	if deviceCreateOutp != nil {
		createdResourceId = &deviceCreateOutp.Uid
	}
	if err != nil {
		return nil, &CreateError{
			CreatedResourceId: createdResourceId,
			Err:               err,
		}
	}

	client.Logger.Println("reading specific device uid")

	asaReadSpecOutp, err := device.ReadSpecific(ctx, client, *device.NewReadSpecificInput(
		deviceCreateOutp.Uid,
	))
	if err != nil {
		return nil, &CreateError{
			CreatedResourceId: createdResourceId,
			Err:               err,
		}
	}

	client.Logger.Println("waiting for asa config state done")

	// poll until asa config state done
	err = retry.Do(asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid), *retry.NewOptionsWithLogger(client.Logger))

	// error during polling, but we maybe able to handle it
	if err != nil {
		// no idea what I am doing here, but that is what the cdo frontend ui is doing
		if createInp.IgnoreCertifcate {
			// update device with ignore certificate
			client.Logger.Println("retrying with ignore certificate")
			_, err := device.Update(ctx, client, device.UpdateInput{
				Uid:              deviceCreateOutp.Uid,
				IgnoreCertifcate: true,
			})
			if err != nil {
				return nil,
					&CreateError{
						CreatedResourceId: createdResourceId,
						Err:               fmt.Errorf("error while updating config to ignore certificate, cause=%w", err),
					}

			}
		} else {
			return nil, &CreateError{
				CreatedResourceId: createdResourceId,
				Err:               err,
			}
		}
	}

	// encrypt credentials for SDC on prem lar
	var publicKey *sdc.PublicKey
	if strings.EqualFold(deviceCreateOutp.LarType, "SDC") {

		// on-prem lar requires encryption
		client.Logger.Println("decrypting public key from sdc for encrpytion")

		if deviceCreateOutp.LarUid == "" {
			return nil, &CreateError{
				CreatedResourceId: createdResourceId,
				Err:               fmt.Errorf("sdc uid not found"),
			}

		}

		// read lar public key
		larReadRes, err := sdc.ReadByUid(ctx, client, sdc.ReadByUidInput{
			SdcUid: deviceCreateOutp.LarUid,
		})
		if err != nil {
			return nil, &CreateError{
				CreatedResourceId: createdResourceId,
				Err:               err,
			}
		}
		publicKey = &larReadRes.PublicKey
	}

	// update asa config credentials
	client.Logger.Println("updating asa config credentials")

	asaConfigUpdateInp := asaconfig.NewUpdateInput(
		asaReadSpecOutp.SpecificUid,
		createInp.Username,
		createInp.Password,
		publicKey,
		asaReadSpecOutp.State,
	)
	_, err = asaconfig.Update(ctx, client, *asaConfigUpdateInp)
	if err != nil {
		return nil, &CreateError{
			CreatedResourceId: createdResourceId,
			Err:               err,
		}
	}

	// poll until asa config state done
	client.Logger.Println("waiting for device to reach state done")

	err = retry.Do(asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid), *retry.NewOptionsWithLogger(client.Logger))
	if err != nil {
		return nil, &CreateError{
			CreatedResourceId: createdResourceId,
			Err:               err,
		}
	}

	// successful
	createOutp := CreateOutput{
		Uid:           deviceCreateOutp.Uid,
		Name:          deviceCreateOutp.Name,
		DeviceType:    deviceCreateOutp.DeviceType,
		Host:          deviceCreateOutp.Host,
		Port:          deviceCreateOutp.Port,
		SocketAddress: deviceCreateOutp.SocketAddress,
		SdcUid:        deviceCreateOutp.LarUid,
		SdcType:       deviceCreateOutp.LarType,
	}
	return &createOutp, nil
}
