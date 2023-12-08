package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func HashObject(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
