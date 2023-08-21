package model

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	KeyId    string `json:"keyId,omitempty"`
}

func NewEncryptedCredentials(encryptedUsername, encryptedPassword, keyId string) Credentials {
	return Credentials{
		Username: encryptedUsername,
		Password: encryptedPassword,
		KeyId:    keyId,
	}
}

func NewCredentials(username, password string) Credentials {
	return Credentials{
		Username: username,
		Password: password,
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
