package portal

import (
	"bytes"
	"net/url"
	"os"
	"time"
)

type Config struct {
	PushInterval     time.Duration
	LogInterval      time.Duration
	AutoPushDelay    time.Duration
	TrustLevel       TrustLevel
	AuthUsername     string
	AuthPassword     string
	AuthPasswordFile string
}

func (c *Config) ParseCredentials() (*url.Userinfo, error) {
	if c.AuthUsername == "" {
		return nil, nil
	}

	p := c.AuthPassword

	if c.AuthPasswordFile != "" {
		b, err := os.ReadFile(c.AuthPasswordFile)
		if err != nil {
			return nil, err
		}

		p = string(bytes.TrimSpace(b))
	}

	return url.UserPassword(c.AuthUsername, p), nil
}
