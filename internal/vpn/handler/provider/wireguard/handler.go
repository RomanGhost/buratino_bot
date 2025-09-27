package wireguard

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider/data"
)

type Address struct {
	IPv4 string `json:"ipv4Address"`
	IPv6 string `json:"ipv6Address"`
}

// Структуры для API ответов
type Client struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Enabled      bool   `json:"enabled"`
	PublicKey    string `json:"publicKey"`
	PresharedKey string `json:"presharedKey"`
	PrivateKey   string `json:""`
	Address
	DNS                 []string `json:"dns"`
	MTU                 int      `json:"mtu"`
	AllowedIPs          []string `json:"allowedIps"`
	PersistentKeepalive int      `json:"persistentKeepalive"`
	Endpoint            *string  `json:"endpoint"`
}

type ClientsResponse struct {
	Clients []Client
}

// Клиент для работы с wg-easy API
type WgEasyClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// Создание нового клиента
func NewWgEasyClient(apiURL, username, password string) *WgEasyClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &WgEasyClient{
		baseURL:    apiURL,
		username:   username,
		password:   password,
		httpClient: client,
	}
}

// Выполнение HTTP запроса
func (c *WgEasyClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("ошибка маршалинга JSON: %v", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.SetBasicAuth(c.username, c.password)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return resp, nil
}

// Получение списка всех клиентов
func (c *WgEasyClient) GetClients() ([]Client, error) {
	resp, err := c.makeRequest("GET", "/api/client", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	// data, _ := io.ReadAll(resp.Body)
	// log.Printf("[DEBUG] Response data: %s", data)
	var clientsResp []Client
	if err := json.NewDecoder(resp.Body).Decode(&clientsResp); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return clientsResp, nil
}

func (c *WgEasyClient) GetClientsByName(name string) ([]Client, error) {
	clients, err := c.GetClients()
	if err != nil {
		return nil, err
	}

	resultClients := make([]Client, 0)
	for _, client := range clients {
		if client.Name == name {
			resultClients = append(resultClients, client)
		}
	}
	return resultClients, nil
}

// Создание нового клиента
func (c *WgEasyClient) CreateClient(name string) error {
	body := map[string]any{"name": name, "expiresAt": nil}
	resp, err := c.makeRequest("POST", "/api/client", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

func (c *WgEasyClient) CreateKey(name string) (*data.KeyConnectData, error) {
	createClientName := fmt.Sprintf("%s-%d", name, time.Now().UTC().Unix())
	err := c.CreateClient(createClientName)
	if err != nil {
		return nil, fmt.Errorf("error client create: %v", err)
	}

	clientsNew, err := c.GetClientsByName(createClientName)
	if err != nil {
		return nil, fmt.Errorf("server error get client: %v", err)
	}
	if len(clientsNew) < 1 {
		return nil, fmt.Errorf("error get new client")
	}
	newClient := clientsNew[0]

	connectContent, err := c.GetConfigurationClientById(newClient.ID)
	if err != nil {
		return nil, fmt.Errorf("error create key: %v", err)
	}

	return &data.KeyConnectData{
		ID:          newClient.ID,
		Name:        createClientName,
		ConnectData: connectContent,
	}, nil

}

func (c *WgEasyClient) GetConfigurationClientById(id int) (string, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/client/%d/configuration", id), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 && resp.StatusCode > 300 {
		return "", fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)

	return string(data), nil

}

// Удаление клиента
func (c *WgEasyClient) DeleteAccessKey(keyID int) error {
	endpoint := fmt.Sprintf("/api/client/%d", keyID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}
