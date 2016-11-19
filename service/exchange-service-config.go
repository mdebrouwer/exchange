package service

type ExchangeServiceConfig struct {
	Host       string     `json:"Host"`
	Port       int32      `json:"Port"`
	AuthConfig AuthConfig `json:"AuthConfig"`
}

type AuthConfig struct {
	CookieSigningKey    string `json:"CookieSigningKey"`
	CookieEncryptionKey string `json:"CookieEncryptionKey"`
	BoltPath            string `json:"BoltPath"`
}
