package web_utils

import (
	"encoding/base64"
	"fmt"
	"net/url"
)

func EncodeRedirectState(state []byte) string {
	return base64.URLEncoding.EncodeToString(state)
}

func DecodeRedirectState(state string) (string, error) {
	decodeString, err := base64.URLEncoding.DecodeString(state)
	return string(decodeString), err
}

func EncodeUrlAndQuery(url *url.URL) string {
	full := fmt.Sprintf("%s?%s", url.Path, url.RawQuery)
	return EncodeRedirectState([]byte(full))
}
