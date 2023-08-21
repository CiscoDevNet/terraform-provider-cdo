package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	internalRsa "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/crypto/internal/rsa"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model"
)

func EncryptCredentials(key model.PublicKey, username, password string) (model.Credentials, error) {
	ciper, err := internalRsa.NewCiper(key.EncodedKey)
	if err != nil {
		return model.Credentials{}, err
	}
	encryptedUsername, err := ciper.Encrypt(username)
	if err != nil {
		return model.Credentials{}, err
	}
	encryptedPassword, err := ciper.Encrypt(password)
	if err != nil {
		return model.Credentials{}, err
	}

	return model.NewEncryptedCredentials(encryptedUsername, encryptedPassword, key.KeyId), nil
}

func DecryptBase64EncodedPkcs1v15Value(privateKey *rsa.PrivateKey, cipherText []byte) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(string(cipherText))
	if err != nil {
		return "", err
	}

	secret, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func MustDecryptBase64EncodedPkcs1v15Value(privateKey *rsa.PrivateKey, cipherText []byte) string {
	secret, err := DecryptBase64EncodedPkcs1v15Value(privateKey, cipherText)
	if err != nil {
		panic(err)
	}

	return secret
}

func Base64PublicKeyFromRsaKey(key *rsa.PrivateKey) (string, error) {
	pub := key.Public()

	bytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}

	pubPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		},
	)

	base64Bytes := make([]byte, base64.StdEncoding.EncodedLen(len(pubPem)))
	base64.StdEncoding.Encode(base64Bytes, pubPem)

	return string(base64Bytes), nil
}

func MustBase64PublicKeyFromRsaKey(key *rsa.PrivateKey) string {
	pubKey, err := Base64PublicKeyFromRsaKey(key)
	if err != nil {
		panic(err)
	}

	return pubKey
}
