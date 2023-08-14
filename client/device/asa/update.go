package asa

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid      string `json:"-"`
	Name     string `json:"name"`
	Location string
	Username string
	Password string
}

type UpdateOutput = device.UpdateOutput

func NewUpdateInput(uid string, name string, username string, password string) *UpdateInput {
	return &UpdateInput{
		Uid:      uid,
		Name:     name,
		Username: username,
		Password: password,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating asa device")

	if isSpecificDeviceIsRequired(updateInp) {

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

		if updateInp.Username != "" || updateInp.Password != "" {
			var publicKey *sdc.PublicKey
			if strings.EqualFold(asaReadOutp.LarType, "SDC") {
				if asaReadOutp.LarUid == "" {
					return nil, fmt.Errorf("sdc uid not found")
				}

				larReadRes, err := sdc.ReadByUid(ctx, client, sdc.ReadByUidInput{
					SdcUid: asaReadOutp.LarUid,
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

			if err := retry.Do(asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid), retry.DefaultOpts); err != nil {
				return nil, err
			}
		}

		if updateInp.Location != "" {
			_, err := asaconfig.UpdateLocation(ctx, client, asaconfig.UpdateLocationOptions{
				SpecificUid: asaReadSpecOutp.SpecificUid,
				Location:    updateInp.Location,
			})
			if err != nil {
				return nil, err
			}

			if err := retry.Do(asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid), retry.DefaultOpts); err != nil {
				return nil, err
			}
		}
	}

	url := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	req := client.NewPut(ctx, url, updateInp)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	if err := retry.Do(UntilStateDoneAndConnectivityOk(ctx, client, outp.Uid), retry.DefaultOpts); err != nil {
		return nil, err
	}

	return &outp, nil
}

func isSpecificDeviceIsRequired(updateInput UpdateInput) bool {
	return updateInput.Username != "" || updateInput.Password != "" || updateInput.Location != ""
}
