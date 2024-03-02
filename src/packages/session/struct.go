package session

import (
	"WTN-Sniper/src/packages/proxy_manager"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type Session struct {
	HttpClient httpClient
	*proxy_manager.Proxy_Manager
}

type httpClient struct {
	TLS_client    tls_client.HttpClient
	TimeOut       int
	ClientProfile profiles.ClientProfile
}
