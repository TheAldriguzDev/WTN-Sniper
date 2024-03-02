package proxy_manager

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func New(proxy_file string) *Proxy_Manager {
	proxyList, err := ReadProxyFile(proxy_file)
	if err != nil {
		panic(err)
	}

	return &Proxy_Manager{
		Proxy_file: proxy_file,
		Proxy_list: proxyList,
	}
}

func ReadProxyFile(file_name string) ([]string, error) {
	if file_name == "localhost" {
		return []string{fmt.Sprintf("127.0.0.1:8888")}, nil
	}
	file_content, err := ioutil.ReadFile(fmt.Sprintf("./Proxies/%s", file_name))

	if err != nil {
		var emptyArray []string
		return emptyArray, err
	}

	proxyList := strings.Split(string(file_content), "\n")
	return proxyList, nil
}

func RandomPort() int {
	// Inizializza il generatore di numeri casuali con un seme univoco
	rand.Seed(time.Now().UnixNano())

	// Genera un numero casuale nell'intervallo specificato
	return rand.Intn(65536-49152) + 49152
}

func (p *Proxy_Manager) GetRandomProxy(proxyList []string) (string, error) {
	if len(proxyList) <= 0 {
		return "", fmt.Errorf("Proxy list is exmpty!")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(proxyList))
	return proxyList[randomIndex], nil
}

func (p *Proxy_Manager) FormatProxy(proxy string) (string, error) {
	proxy = strings.ReplaceAll(proxy, "\r", "")
	proxy_data := strings.Split(strings.TrimSpace(proxy), ":")

	if len(proxy_data) == 2 {
		// IP PORT proxy (no auth)
		proxyUrl := fmt.Sprintf("http://%s:%s", proxy_data[0], proxy_data[1])
		p.Current_Proxy = proxyUrl
		return proxyUrl, nil
	} else if len(proxy_data) == 4 {
		proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxy_data[2], proxy_data[3], proxy_data[0], proxy_data[1])
		p.Current_Proxy = proxyUrl
		return proxyUrl, nil
	} else {
		return "", fmt.Errorf("proxy format not supported")
	}
}
