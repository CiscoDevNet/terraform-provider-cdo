package iosconfig

import (
	"context"
	"encoding/json"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
)

type UpdateInput struct {
	SpecificUid string
	PublicKey   *model.PublicKey

	Username string
	Password string
}

type UpdateOutput struct {
	Uid string `json:"uid"`
}

func NewUpdateInput(specificUid string, username string, password string, publicKey *model.PublicKey) *UpdateInput {
	return &UpdateInput{
		SpecificUid: specificUid,
		Username:    username,
		Password:    password,
		PublicKey:   publicKey,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating iosconfig")

	url := url.UpdateDevice(client.BaseUrl(), updateInp.SpecificUid)

	creds, err := makeCredentials(updateInp)
	if err != nil {
		return nil, err
	}

	req := client.NewPut(ctx, url, makeReqBody(creds))

	var outp UpdateOutput
	err = req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}

func makeReqBody(creds []byte) *UpdateBody {
	return &UpdateBody{
		Credentials: string(creds),
		SmContext: SmContext{
			AcceptCert: true,
		},
	}
}

type UpdateBody struct {
	Credentials string    `json:"credentials"`
	SmContext   SmContext `json:"stateMachineContext"`
}

type SmContext struct {
	AcceptCert bool `json:"acceptCert"`
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
