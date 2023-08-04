package iosconfig

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/go-client/connector/sdc"
	internalRsa "github.com/CiscoDevnet/go-client/internal/crypto/rsa"
	internalHttp "github.com/CiscoDevnet/go-client/internal/http"
	"github.com/CiscoDevnet/go-client/internal/jsonutil"
	"github.com/jarcoal/httpmock"
)

func TestIosConfigUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "unit-test-username"
	password := "not a real password"

	rsaKeyBits := 512
	rsaKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		t.Fatal("could not generate rsa key")
	}

	validIosConfig := ReadOutput{
		Uid:   iosConfigUid,
		State: IosConfigStateDone,
	}

	testCases := []struct {
		testName   string
		input      UpdateInput
		setupFunc  func(input UpdateInput, t *testing.T)
		assertFunc func(output *UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates iOS config",
			input: UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(iosConfigUid),
					func(r *http.Request) (*http.Response, error) {
						requestBody, err := internalHttp.ReadRequestBody[UpdateBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						expectedBody := UpdateBody{
							SmContext: SmContext{
								AcceptCert: true,
							},
							Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
						}

						if !reflect.DeepEqual(expectedBody, *requestBody) {
							t.Errorf("expected request body to equal: %+v, got: %+v", expectedBody, requestBody)
						}

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: iosConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(UpdateOutput{Uid: iosConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validIosConfig, output)
				}
			},
		},

		{
			testName: "successfully updates iOS config when encrypting credentials",
			input: UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
				PublicKey: &sdc.PublicKey{
					KeyId:      "12341234-1234-1234-1234-123412341234",
					Version:    2,
					EncodedKey: internalRsa.MustBase64PublicKeyFromRsaKey(rsaKey),
				},
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(iosConfigUid),
					func(r *http.Request) (*http.Response, error) {
						requestBody, err := internalHttp.ReadRequestBody[UpdateBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						if !requestBody.SmContext.AcceptCert {
							t.Errorf("expected 'SmContext.AcceptCert' to true got: %t", requestBody.SmContext.AcceptCert)
						}

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.Credentials))
						if err != nil {
							t.Fatalf("could not unmarshal credentials because: %s", err.Error())
						}

						decryptedUsername := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Username))
						if input.Username != decryptedUsername {
							t.Errorf(`expected decrypted username to equal '%s', got: '%s'`, input.Username, decryptedUsername)
						}

						decryptedPassword := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Password))
						if input.Password != decryptedPassword {
							t.Errorf(`expected decrypted password to equal '%s', got: '%s'`, input.Password, decryptedPassword)
						}

						if input.PublicKey.KeyId != credentials.KeyId {
							t.Errorf("expected keyId to equal '%s', got: '%s'", input.PublicKey.KeyId, credentials.KeyId)
						}

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: iosConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(UpdateOutput{Uid: iosConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validIosConfig, output)
				}
			},
		},

		{
			testName: "returns error when updating iOS config that does not exist",
			input: UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(input.SpecificUid),
					httpmock.NewStringResponder(404, ""),
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", *output)
				}

				if err == nil {
					t.Fatal("err is nil!")
				}
			},
		},

		{
			testName: "returns error when remote service updating iOS config experiences an issue",
			input: UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(input.SpecificUid),
					httpmock.NewStringResponder(500, ""),
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if output != nil {
					t.Errorf("expected output to be nil, got: %+v", *output)
				}

				if err == nil {
					t.Fatal("err is nil!")
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input, t)

			output, err := Update(context.Background(), *internalHttp.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}
