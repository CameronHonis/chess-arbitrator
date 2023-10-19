package server

import "encoding/json"

type AuthMessageContent struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func (msg *AuthMessageContent) toJsonBytes() ([]byte, error) {
	msgJson, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return msgJson, nil
}
