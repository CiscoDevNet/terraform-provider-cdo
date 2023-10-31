package asa

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/asa/asaconfig"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/retry"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	Uid      string    `json:"-"`
	Name     string    `json:"name"`
	Location string    `json:"-"`
	Username string    `json:"-"`
	Password string    `json:"-"`
	Tags     tags.Type `json:"tags"`
}

type UpdateOutput = device.UpdateOutput

func NewUpdateInput(uid string, name string, username string, password string, tags tags.Type) *UpdateInput {
	return &UpdateInput{
		Uid:      uid,
		Name:     name,
		Username: username,
		Password: password,
		Tags:     tags,
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

		asaReadOutp, err := device.ReadByUid(ctx, client, *device.NewReadByUidInput(
			updateInp.Uid,
		))
		if err != nil {
			return nil, err
		}

		if updateInp.Username != "" || updateInp.Password != "" {
			var publicKey *model.PublicKey
			if strings.EqualFold(asaReadOutp.ConnectorType, "SDC") {
				if asaReadOutp.ConnectorUid == "" {
					return nil, fmt.Errorf("connector uid not found")
				}

				connectorReadRes, err := connector.ReadByUid(ctx, client, connector.ReadByUidInput{
					ConnectorUid: asaReadOutp.ConnectorUid,
				})
				if err != nil {
					return nil, err
				}
				publicKey = &connectorReadRes.PublicKey
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

			if err := retry.Do(
				ctx,
				asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid),
				retry.NewOptionsBuilder().
					Title("waiting for update ASA credentials to be done").
					Retries(retry.DefaultRetries).
					Delay(retry.DefaultDelay).
					Timeout(retry.DefaultTimeout).
					EarlyExitOnError(true).
					Build(),
			); err != nil {
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

			if err := retry.Do(
				ctx,
				asaconfig.UntilStateDone(ctx, client, asaReadSpecOutp.SpecificUid),
				retry.NewOptionsBuilder().
					Title("waiting for update ASA location to be done").
					Retries(retry.DefaultRetries).
					Delay(retry.DefaultDelay).
					Timeout(retry.DefaultTimeout).
					EarlyExitOnError(true).
					Build(),
			); err != nil {
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

	if err := retry.Do(
		ctx,
		UntilStateDoneAndConnectivityOk(ctx, client, outp.Uid),
		retry.NewOptionsBuilder().
			Title("waiting for ASA to be done with connectivity status OK").
			Retries(retry.DefaultRetries).
			Delay(retry.DefaultDelay).
			Timeout(retry.DefaultTimeout).
			EarlyExitOnError(true).
			Build(),
	); err != nil {
		return nil, err
	}

	return &outp, nil
}

func isSpecificDeviceIsRequired(updateInput UpdateInput) bool {
	return updateInput.Username != "" || updateInput.Password != "" || updateInput.Location != ""
}
