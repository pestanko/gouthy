package web_utils

import (
	"encoding/base64"
	"encoding/json"
)

func JsonEncode(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return ToBase64Url(bytes), nil
}

func JsonDecode(data string, result interface{}) error {
	bytes, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, result); err != nil {
		return err
	}

	return nil
}

func ToBase64Url(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}
