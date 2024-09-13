package auth

import (
	"context"
	"encoding/base64"
)

type BasicAuth struct {
	Password string
	Username string
}

func (b BasicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	auth := b.Username + ":" + b.Password
	enc := base64.StdEncoding.EncodeToString([]byte(auth))

	return map[string]string{
		"authorization": "Basic " + enc,
	}, nil
}

func (b BasicAuth) RequireTransportSecurity() bool {
	return true
}
