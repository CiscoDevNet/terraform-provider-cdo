package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type RsaCiper struct {
	pub *rsa.PublicKey
}

func NewCiper(base64EncodedPublicKey string) (*RsaCiper, error) {
	encodedKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(encodedKey))
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &RsaCiper{
		pub.(*rsa.PublicKey),
	}, nil
}

func (ciper *RsaCiper) Encrypt(msg string) (string, error) {
	rsaEncoded, err := rsa.EncryptPKCS1v15(rand.Reader, ciper.pub, []byte(msg))
	if err != nil {
		return "", err
	}
	base64RsaEncoded := base64.StdEncoding.EncodeToString(rsaEncoded)

	return base64RsaEncoded, nil
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
