package security

import (
	"builder/model"
	"encoding/base64"
	"errors"
	"strings"
)

// ParseBasicToken returns parsed basic token
func ParseBasicToken(token string) (*model.BasicToken, error) {
	splitted := strings.Split(token, " ")

	if len(splitted) < 2 || strings.ToLower(splitted[0]) != "basic" {
		return nil, errors.New("Not Basic Token")
	}
	decoded, err := base64.StdEncoding.DecodeString(splitted[1])
	if err != nil {
		return nil, errors.New("Token should be Base64 encoded")
	}
	splitted = strings.Split(string(decoded), ":")
	if len(splitted) < 2 {
		return nil, errors.New("Not exists username and password")
	}
	return &model.BasicToken{
		Raw:      token,
		Username: splitted[0],
		Password: splitted[1],
	}, nil
}
