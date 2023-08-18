package genericssh_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/device/genericssh"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto"
	internalHttp "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/jsonutil"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/statemachine/state"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestGenericSshUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.UpdateOutput{
		Uid:   genericSshUid,
		Name:  genericSshName,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully update Generic SSH",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, validGenericSsh),
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "return error when updating Generic SSH error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				genericssh.NewUpdateInput(
					genericSshUid,
					"",
					"",
					"",
					nil,
				),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestGenericSshUpdateWithEncryptedCredentials(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.UpdateOutput{
		Uid:   genericSshUid,
		Name:  genericSshName,
		State: state.DONE,
	}

	// generate public key for this test
	keyId := "unit-test-key-id"
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 512)
	assert.Nil(t, err, "failed to generate rsa key")
	publicKey := model.NewPublicKey(
		crypto.MustBase64PublicKeyFromRsaKey(rsaPrivateKey),
		123,
		keyId,
	)

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully update Generic SSH with encrypted credentials",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					func(r *http.Request) (*http.Response, error) {
						// read credentials from request
						body, err := internalHttp.ReadRequestBody[genericssh.UpdateBody](r)
						assert.Nil(t, err)
						creds, err := jsonutil.UnmarshalStruct[model.Credentials]([]byte(body.Credentials))
						assert.Nil(t, err)

						// validate credentials is indeed correctly encrypted
						assert.Equal(t, creds.KeyId, keyId)
						decryptedUsername := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaPrivateKey, []byte(creds.Username))
						assert.Equal(t, decryptedUsername, genericSshUsername)
						decryptedPassword := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaPrivateKey, []byte(creds.Password))
						assert.Equal(t, decryptedPassword, genericSshPassword)

						return httpmock.NewJsonResponse(http.StatusOK, validGenericSsh)
					},
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "return error when update Generic SSH with encrypted credentials error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					func(r *http.Request) (*http.Response, error) {
						// read credentials from request
						body, err := internalHttp.ReadRequestBody[genericssh.UpdateBody](r)
						assert.Nil(t, err)
						creds, err := jsonutil.UnmarshalStruct[model.Credentials]([]byte(body.Credentials))
						assert.Nil(t, err)

						// validate credentials is indeed correctly encrypted
						assert.Equal(t, creds.KeyId, keyId)
						decryptedUsername := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaPrivateKey, []byte(creds.Username))
						assert.Equal(t, decryptedUsername, genericSshUsername)
						decryptedPassword := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaPrivateKey, []byte(creds.Password))
						assert.Equal(t, decryptedPassword, genericSshPassword)

						return httpmock.NewJsonResponse(http.StatusInternalServerError, "internal server error")
					},
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				genericssh.NewUpdateInput(
					genericSshUid,
					genericSshConnectorUid,
					genericSshUsername,
					genericSshPassword,
					&publicKey,
				),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}

func TestGenericSshUpdateWithCredentials(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validGenericSsh := genericssh.UpdateOutput{
		Uid:   genericSshUid,
		Name:  genericSshName,
		State: state.DONE,
	}

	testCases := []struct {
		testName   string
		targetUid  string
		setupFunc  func()
		assertFunc func(output *genericssh.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName:  "successfully update Generic SSH with credentials",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					func(r *http.Request) (*http.Response, error) {
						// read credentials from request
						body, err := internalHttp.ReadRequestBody[genericssh.UpdateBody](r)
						assert.Nil(t, err)
						creds, err := jsonutil.UnmarshalStruct[model.Credentials]([]byte(body.Credentials))
						assert.Nil(t, err)

						// validate credentials is indeed correctly encrypted
						assert.Zero(t, creds.KeyId, "key id should have zero value")
						assert.Equal(t, creds.Username, genericSshUsername)
						assert.Equal(t, creds.Password, genericSshPassword)

						return httpmock.NewJsonResponse(http.StatusOK, validGenericSsh)
					},
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validGenericSsh, *output)
			},
		},
		{
			testName:  "return error when update Generic SSH with credentials error",
			targetUid: genericSshUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodPut,
					url.UpdateDevice(baseUrl, genericSshUid),
					func(r *http.Request) (*http.Response, error) {
						// read credentials from request
						body, err := internalHttp.ReadRequestBody[genericssh.UpdateBody](r)
						assert.Nil(t, err)
						creds, err := jsonutil.UnmarshalStruct[model.Credentials]([]byte(body.Credentials))
						assert.Nil(t, err)

						// validate credentials is indeed correctly encrypted
						assert.Zero(t, creds.KeyId, "key id should have zero value")
						assert.Equal(t, creds.Username, genericSshUsername)
						assert.Equal(t, creds.Password, genericSshPassword)

						return httpmock.NewJsonResponse(http.StatusInternalServerError, "internal server error")
					},
				)
			},
			assertFunc: func(output *genericssh.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := genericssh.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				genericssh.NewUpdateInput(
					genericSshUid,
					genericSshConnectorUid,
					genericSshUsername,
					genericSshPassword,
					nil,
				),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
