package asa

import (
	"context"
	"fmt"	
	"github.com/cisco-lockhart/go-client/connector/sdc"
	"github.com/cisco-lockhart/go-client/internal/device/asaconfig"
	"strings"

	"github.com/cisco-lockhart/go-client/device"
	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type UpdateInput struct {
	Uid       string `json:"-"`
	Name      string `json:"name"`
	Username  string
	Password  string
}

type UpdateOutput = device.UpdateOutput

func NewUpdateInput(uid string, name string, username string, password string) *UpdateInput {
	return &UpdateInput{
		Uid:       uid,
		Name:      name,
		Username:  username,
		Password:  password,
	}
}

func NewUpdateRequest(ctx context.Context, client http.Client, updateInp UpdateInput) *http.Request {

	url := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	return req
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating asa device")

	if updateInp.Username != "" {

		asaReadSpecOutp, err := device.ReadSpecific(ctx, client, *device.NewReadSpecificInput(
			updateInp.Uid,
		))
		if err != nil {
			return nil, err
		}

		asaReadOutp, err := device.Read(ctx, client, *device.NewReadInput(
			updateInp.Uid,
		))
		if err != nil {
			return nil, err
		}

		var publicKey *sdc.PublicKey
		if strings.EqualFold(asaReadOutp.LarType, "SDC") {
			if asaReadOutp.LarUid == "" {
				return nil, fmt.Errorf("sdc uid not found")
			}
	
			larReadRes, err := sdc.ReadByUid(ctx, client, sdc.ReadInput{
				LarUid: asaReadOutp.LarUid,
			})
			if err != nil {
				return nil, err
			}
			publicKey = &larReadRes.PublicKey
		}

		updateAsaConfigInp := asaconfig.NewUpdateInput(
			asaReadSpecOutp.SpecificUid,
			updateInp.Username,
			updateInp.Password,
			publicKey,
			asaReadSpecOutp.State,
		)
		_, err = asaconfig.UpdateCredentials(ctx, client, *updateAsaConfigInp)
		if err != nil {
			_ = fmt.Errorf("failed to update credentials for ASA device: %s", err.Error())
			return nil, err
		}
	}

	req := NewUpdateRequest(ctx, client, updateInp)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
