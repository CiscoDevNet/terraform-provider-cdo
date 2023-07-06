package ios

import (
	"context"
	"fmt"
	"github.com/cisco-lockhart/go-client/device/ios/iosconfig"
	"strings"

	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/device/sdc"
	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/retry"
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

	specificUid string `json:"-"`
}

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

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating ios device")

	deviceCreateOutp, err := device.Create(ctx, client, *device.NewCreateRequestInput(
		createInp.Name, "IOS", createInp.LarUid, createInp.LarType, createInp.Ipv4, false, createInp.IgnoreCertifcate,
	))
	if err != nil {
		return nil, err
	}

	client.Logger.Println("reading specific device uid")

	iosReadSpecOutp, err := device.ReadSpecific(ctx, client, *device.NewReadSpecificInput(
		deviceCreateOutp.Uid,
	))
	if err != nil {
		return nil, err
	}

	client.Logger.Println("waiting for ios config state done")

	// poll until ios config state done
	err = retry.Do(iosconfig.UntilStateDone(ctx, client, iosReadSpecOutp.SpecificUid), *retry.NewOptionsWithLogger(client.Logger))

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
				return nil, fmt.Errorf("error while updating config to ignore certificate, cause=%w", err)
			}
		} else {
			return nil, err
		}
	}

	// encrypt credentials for SDC on prem lar
	var publicKey *sdc.PublicKey
	if strings.EqualFold(deviceCreateOutp.LarType, "SDC") {

		// on-prem lar requires encryption
		client.Logger.Println("decrypting public key from sdc for encrpytion")

		if deviceCreateOutp.LarUid == "" {
			return nil, fmt.Errorf("sdc uid not found")
		}

		// read lar public key
		larReadRes, err := sdc.Read(ctx, client, sdc.ReadInput{
			LarUid: deviceCreateOutp.LarUid,
		})
		if err != nil {
			return nil, err
		}
		publicKey = &larReadRes.PublicKey
	}

	// update ios config credentials
	client.Logger.Println("updating ios config credentials")

	iosConfigUpdateInp := iosconfig.NewUpdateInput(
		iosReadSpecOutp.SpecificUid,
		createInp.Username,
		createInp.Password,
		publicKey,
	)
	_, err = iosconfig.Update(ctx, client, *iosConfigUpdateInp)
	if err != nil {
		return nil, err
	}

	// poll until ios config state done
	client.Logger.Println("waiting for device to reach state done")

	err = retry.Do(iosconfig.UntilStateDone(ctx, client, iosReadSpecOutp.SpecificUid), *retry.NewOptionsWithLogger(client.Logger))
	if err != nil {
		return nil, err
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

		specificUid: iosReadSpecOutp.SpecificUid,
	}
	return &createOutp, nil
}
