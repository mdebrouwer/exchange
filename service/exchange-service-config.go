package service

import (
	"bufio"
	"os"
	"strings"
)

type ExchangeServiceConfig struct {
	Host       string     `json:"Host"`
	Port       int32      `json:"Port"`
	AuthConfig AuthConfig `json:"AuthConfig"`
}

type AuthConfig struct {
	CookieSigningKey    string `json:"CookieSigningKey"`
	CookieEncryptionKey string `json:"CookieEncryptionKey"`
	BoltPath            string `json:"BoltPath"`
	TokenFile           string `json:"TokenFile"`
}

func (c AuthConfig) Tokens() (values []string, err error) {
	values = make([]string, 0)
	set := make(map[string]bool)
	file, err := os.Open(c.TokenFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if value != "" {
			set[value] = true
		}
	}

	err = scanner.Err()
	if err != nil {
		return
	}

	for k, _ := range set {
		values = append(values, k)
	}

	return
}
