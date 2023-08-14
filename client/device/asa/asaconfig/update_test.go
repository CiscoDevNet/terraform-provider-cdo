package asaconfig

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	h "net/http"
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
	assert.Nil(t, err, "could not generate rsa key")

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
						assert.Nil(t, err)

						expectedBody := updateBody{
							State:       "CERT_VALIDATED",
							Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
						}
						assert.Equal(t, expectedBody, *requestBody)

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
						assert.Nil(t, err)
						assert.Equal(t, requestBody.State, "CERT_VALIDATED", fmt.Sprintf("expected 'State' to equal 'CERT_VALIDATED', got: %s", requestBody.State))

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.Credentials))
						assert.Nil(t, err)

						decryptedUsername := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Username))
						assert.Equal(t, input.Username, decryptedUsername, fmt.Sprintf(`expected decrypted username to equal '%s', got: '%s'`, input.Username, decryptedUsername))

						decryptedPassword := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Password))
						assert.Equal(t, input.Password, decryptedPassword, fmt.Sprintf(`expected decrypted password to equal '%s', got: '%s'`, input.Password, decryptedPassword))
						assert.Equal(t, input.PublicKey.KeyId, credentials.KeyId, fmt.Sprintf("expected keyId to equal '%s', got: '%s'", input.PublicKey.KeyId, credentials.KeyId))

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
				assert.Nil(t, output)
				assert.NotNil(t, err)
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
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input, t)

			output, err := Update(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

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
						assert.Nil(t, err)

						expectedBody := updateCredentialsBodyWithState{
							State: "WAIT_FOR_USER_TO_UPDATE_CREDS",
							SmContext: SmContext{
								Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
							},
						}
						assert.Equal(t, expectedBody, *requestBody)

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
						assert.Nil(t, err)

						expectedState := "WAIT_FOR_USER_TO_UPDATE_CREDS"
						assert.Equal(t, requestBody.State, expectedState)

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.SmContext.Credentials))
						assert.Nil(t, err)

						decryptedUsername := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Username))
						assert.Equal(t, input.Username, decryptedUsername)

						decryptedPassword := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Password))
						assert.Equal(t, input.Password, decryptedPassword)
						assert.Equal(t, input.PublicKey.KeyId, credentials.KeyId)

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					})
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
						assert.Nil(t, err)

						expectedBody := updateCredentialsBody{
							SmContext: SmContext{
								Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
							},
						}
						assert.Equal(t, expectedBody, *requestBody)

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					},
				)
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
						assert.Nil(t, err)

						credentials, err := jsonutil.UnmarshalStruct[credentials]([]byte(requestBody.SmContext.Credentials))
						assert.Nil(t, err)

						decryptedUsername := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Username))
						assert.Equal(t, input.Username, decryptedUsername)

						decryptedPassword := internalRsa.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Password))
						assert.Equal(t, input.Password, decryptedPassword)
						assert.Equal(t, input.PublicKey.KeyId, credentials.KeyId)

						return httpmock.NewJsonResponse(200, UpdateOutput{Uid: asaConfigUid})
					})
			},

			assertFunc: func(output *UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, UpdateOutput{Uid: asaConfigUid}, *output)
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
				assert.Nil(t, output)
				assert.NotNil(t, err)
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
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input, t)

			output, err := UpdateCredentials(context.Background(), *http.MustNewWithDefault("https://unittest.cdo.cisco.com", "a_valid_token"), testCase.input)

			testCase.assertFunc(output, err, t)
		})
	}
}

func buildUpdateAsaConfigUrl(uid string) string {
	return fmt.Sprintf("/aegis/rest/v1/services/asa/configs/%s", asaConfigUid)
}
