package session

import (
	"WTN-Sniper/src/packages/proxy_manager"
	"fmt"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

/*
La funzione New crea una nuova session. Vengono inizializzati oggetto standard, come la TLS.
Questa struttura, possendendo tutte le informazioni "anagrafiche" del monitor, potrà
probabilmente implentare le funzioni che si interfacciano con il DB (es. refreshare le task)
*/
func New(timeout int, client_profile profiles.ClientProfile, proxy_session *proxy_manager.Proxy_Manager) (*Session, error) {
	random_proxy, err := proxy_session.GetRandomProxy(proxy_session.Proxy_list)
	if err != nil {
		// Add logger logic
		fmt.Println(fmt.Sprintf("An error occured while getting a random proxy: %e", err))
	}

	formatted_proxy, err := proxy_session.FormatProxy(random_proxy)
	if err != nil {
		// Add logger logic
		fmt.Println(fmt.Sprintf("An error occured while formatting the proxy: %e", err))
	}

	tls_session, err := createTlsClient(formatted_proxy, timeout, client_profile)
	if err != nil {
		// Add logger logic
		fmt.Println(fmt.Sprintf("An error occured while creating a tls client: %e", err))
	}

	scraperSession := &Session{
		HttpClient:    *tls_session,
		Proxy_Manager: proxy_session,
	}

	return scraperSession, nil
}

func createTlsClient(proxy_url string, timeout int, client_profile profiles.ClientProfile) (*httpClient, error) {
	jarOptions := []tls_client.CookieJarOption{}
	jar := tls_client.NewCookieJar(jarOptions...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(timeout),
		tls_client.WithClientProfile(client_profile), // If this value is empty it will use the latest chrome profile available
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
		//tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithProxyUrl(proxy_url),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	httpSession := &httpClient{
		TLS_client:    client,
		TimeOut:       timeout,
		ClientProfile: client_profile,
	}

	return httpSession, nil
}

func (s *Session) RotateProxy() (string, error) {
	random_proxy, err := s.Proxy_Manager.GetRandomProxy(s.Proxy_Manager.Proxy_list)
	if err != nil {
		return "", err
	}

	formatted_proxy, err := s.Proxy_Manager.FormatProxy(random_proxy)
	if err != nil {
		return "", err
	}

	err = s.HttpClient.TLS_client.SetProxy(formatted_proxy)
	if err != nil {
		return "", err
	}

	s.Proxy_Manager.Current_Proxy = formatted_proxy
	return formatted_proxy, nil
}
