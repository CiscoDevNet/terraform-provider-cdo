package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

type RsaCiper struct {
	pub *rsa.PublicKey
}

func NewCiper(base64EncodedPublicKey string) (*RsaCiper, error) {
	encodedKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(encodedKey)
	if block == nil {
		return nil, fmt.Errorf("PEM formatted block (i.e. public key) not found")
	}
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
