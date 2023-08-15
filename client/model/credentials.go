package model

type EncryptedCredentials struct {
	EncryptedUsername string `json:"username"`
	EncryptedPassword string `json:"password"`
	KeyId             string `json:"keyId,omitempty"`
}

func NewEncryptedCredentials(encryptedUsername, encryptedPassword, keyId string) EncryptedCredentials {
	return EncryptedCredentials{
		EncryptedUsername: encryptedUsername,
		EncryptedPassword: encryptedPassword,
		KeyId:             keyId,
	}
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewCredentials(encryptedUsername, encryptedPassword string) Credentials {
	return Credentials{
		Username: encryptedUsername,
		Password: encryptedPassword,
	}
}

type PublicKey struct {
	EncodedKey string `json:"encodedKey"`
	Version    int64  `json:"version"`
	KeyId      string `json:"keyId"`
}

func NewPublicKey(encodedKey string, version int64, keyId string) PublicKey {
	return PublicKey{
		EncodedKey: encodedKey,
		Version:    version,
		KeyId:      keyId,
	}
}
