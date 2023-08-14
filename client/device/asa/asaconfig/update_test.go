package asaconfig

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	h "net/http"
	"reflect"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/connector/sdc"
	internalRsa "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto/rsa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/jsonutil"
	"github.com/jarcoal/httpmock"
)

func TestAsaConfigUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "unit-test-username"
	password := "not a real password"

	rsaKeyBits := 512
	rsaKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		t.Fatal("could not generate rsa key")
	}

	validAsaConfig := ReadOutput{
		Uid:   asaConfigUid,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		input      UpdateInput
		setupFunc  func(input UpdateInput, t *testing.T)
		assertFunc func(output *UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates ASA config",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						expectedBody := updateBody{
							State:       "CERT_VALIDATED",
							Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
						}

						if !reflect.DeepEqual(expectedBody, *requestBody) {
							t.Errorf("expected request body to equal: %+v, got: %+v", expectedBody, requestBody)
						}

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
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

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},
		{
			testName: "successfully updates ASA config when encrypting credentials",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
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
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						if requestBody.State != "CERT_VALIDATED" {
							t.Errorf("expected 'State' to equal 'CERT_VALIDATED', got: %s", requestBody.State)
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

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
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

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},
		{
			testName: "returns error when updating ASA config that does not exist",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(input.SpecificUid),
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
			testName: "returns error when remote service updating ASA config experiences an issue",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(input.SpecificUid),
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

			output, err := Update(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestAsaConfigUpdateCredentials(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "unit-test-username"
	password := "not a real password"

	rsaKeyBits := 512
	rsaKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	if err != nil {
		t.Fatal("could not generate rsa key")
	}

	validAsaConfig := ReadOutput{
		Uid:   asaConfigUid,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		input      UpdateInput
		setupFunc  func(input UpdateInput, t *testing.T)
		assertFunc func(output *UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates ASA config credentials",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateCredentialsBodyWithState](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						expectedBody := updateCredentialsBodyWithState{
							State: "WAIT_FOR_USER_TO_UPDATE_CREDS",
							SmContext: SmContext{
								Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
							},
						}

						if !reflect.DeepEqual(expectedBody, *requestBody) {
							t.Errorf("expected request body to equal: %+v, got: %+v", expectedBody, requestBody)
						}

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
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

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},

		{
			testName: "successfully updates ASA config credentials when encrypting credentials",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
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
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateCredentialsBodyWithState](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						expectedState := "WAIT_FOR_USER_TO_UPDATE_CREDS"
						if requestBody.State != expectedState {
							t.Errorf("expected 'State' to equal '%s', got: %s", expectedState, requestBody.State)
						}

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.SmContext.Credentials))
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

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					})
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},

		{
			testName: "successfully updates ASA config credentials with flag to wait for user to update credentials",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateCredentialsBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						expectedBody := updateCredentialsBody{
							SmContext: SmContext{
								Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
							},
						}

						if !reflect.DeepEqual(expectedBody, *requestBody) {
							t.Errorf("expected request body to equal: %+v, got: %+v", expectedBody, requestBody)
						}

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
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

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},

		{
			testName: "successfully updates ASA config credentials when encrypting credentials with flag to wait for user to update credentials",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
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
					buildUpdateAsaConfigUrl(asaConfigUid),
					func(r *h.Request) (*h.Response, error) {
						requestBody, err := http.ReadRequestBody[updateCredentialsBody](r)
						if err != nil {
							t.Fatalf("could not read body because: %s", err.Error())
						}

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.SmContext.Credentials))
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

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					})
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				if err != nil {
					t.Errorf("unexpected error: %s", err.Error())
				}

				if output == nil {
					t.Fatal("output is nil!")
				}

				if !reflect.DeepEqual(UpdateOutput{Uid: asaConfigUid}, *output) {
					t.Errorf("expected: %+v\ngot: %+v", validAsaConfig, output)
				}
			},
		},

		{
			testName: "returns error when updating ASA config credentials that does not exist",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(input.SpecificUid),
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
			testName: "returns error when remote service updating ASA config credentials experiences an issue",
			input: UpdateInput{
				SpecificUid: asaConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildUpdateAsaConfigUrl(input.SpecificUid),
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

			output, err := UpdateCredentials(context.Background(), *http.NewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}

func buildUpdateAsaConfigUrl(uid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/asa/configs/%s", asaConfigUid)
}
