package genericssh

import (
	"context"
	"encoding/json"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
)

type UpdateInput struct {
	Uid       string
	Name      string
	Username  string
	Password  string
	PublicKey *model.PublicKey
}

type UpdateOutput = device.ReadOutput

func NewUpdateInput(uid, name, username, password string, publicKey *model.PublicKey) UpdateInput {
	return UpdateInput{
		Uid:       uid,
		Name:      name,
		Username:  username,
		Password:  password,
		PublicKey: publicKey,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating generic ssh")

	updateUrl := url.UpdateDevice(client.BaseUrl(), updateInp.Uid)

	updateBody, err := makeUpdateBody(updateInp)
	if err != nil {
		return nil, err
	}

	req := client.NewPut(ctx, updateUrl, updateBody)

	var outp UpdateOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

type UpdateBody struct {
	Name        string    `json:"name,omitempty"`
	Credentials string    `json:"credentials,omitempty"`
	SmContext   SmContext `json:"stateMachineContext,omitempty"`
}

type SmContext struct {
	AcceptCert bool `json:"acceptCert"`
}

func makeUpdateBody(updateInp UpdateInput) (UpdateBody, error) {
	updateBody := UpdateBody{
		Name: updateInp.Name,
	}
	if updateInp.Username != "" {
		creds, err := makeCredentials(updateInp)
		if err != nil {
			return updateBody, err
		}
		updateBody.Credentials = string(creds)
	}

	return updateBody, nil
}

func makeCredentials(updateInp UpdateInput) ([]byte, error) {
	if updateInp.PublicKey != nil {
		encryptedCredentials, err := crypto.EncryptCredentials(*updateInp.PublicKey, updateInp.Username, updateInp.Password)
		if err != nil {
			return nil, err
		}
		return json.Marshal(encryptedCredentials)
	}

	return json.Marshal(model.NewCredentials(updateInp.Username, updateInp.Password))
}
