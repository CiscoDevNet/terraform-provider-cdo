package asaconfig

import (
	"context"
	"encoding/json"

	"github.com/cisco-lockhart/go-client/device/sdc"
	"github.com/cisco-lockhart/go-client/internal/crypto/rsa"
	"github.com/cisco-lockhart/go-client/internal/http"
	"github.com/cisco-lockhart/go-client/internal/url"
)

type UpdateInput struct {
	SpecificUid string
	PublicKey   *sdc.PublicKey

	Username string
	Password string
}

type UpdateOutput struct {
	Uid string `json:"uid"`
}

func NewUpdateInput(specificUid string, username string, password string, publicKey *sdc.PublicKey) *UpdateInput {
	return &UpdateInput{
		SpecificUid: specificUid,
		Username:    username,
		Password:    password,
		PublicKey:   publicKey,
	}
}

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating asaconfig")

	url := url.UpdateAsaConfig(client.BaseUrl(), updateInp.SpecificUid)

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

func makeReqBody(creds []byte) *updateBody {
	return &updateBody{
		State:       "CERT_VALIDATED", // question: should this be hardcoded?
		Credentials: string(creds),
	}
}

type updateBody struct {
	State       string `json:"state"`
	Credentials string `json:"credentials"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	KeyId    string `json:"keyId,omitempty"`
}

func encrypt(req UpdateInput) error {
	ciper, err := rsa.NewCiper(req.PublicKey.EncodedKey)
	if err != nil {
		return err
	}
	req.Username, err = ciper.Encrypt(req.Username)
	if err != nil {
		return err
	}
	req.Password, err = ciper.Encrypt(req.Password)
	if err != nil {
		return err
	}

	return nil
}

func makeCredentials(updateInp UpdateInput) ([]byte, error) {
	var creds []byte
	var err error
	if updateInp.PublicKey != nil {
		encrypt(updateInp)
		creds, err = json.Marshal(credentials{
			Username: updateInp.Username,
			Password: updateInp.Password,
			KeyId:    updateInp.PublicKey.KeyId,
		})
	} else {
		creds, err = json.Marshal(credentials{
			Username: updateInp.Username,
			Password: updateInp.Password,
		})
	}
	return creds, err
}
