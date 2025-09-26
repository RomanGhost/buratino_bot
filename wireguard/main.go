package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Структуры для API ответов
type Client struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	Enabled   bool   `json:"enabled"`
	PublicKey string `json:"publicKey"`
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
func (c *WgEasyClient) CreateClient(name string) (*Client, error) {
	body := map[string]any{"name": name, "expiresAt": nil}
	resp, err := c.makeRequest("POST", "/api/client", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	// data, _ := io.ReadAll(resp.Body)
	// log.Printf("[DEBUG] Response data create client: %s", data)
	var client Client
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &client, nil
}

// Удаление клиента
func (c *WgEasyClient) DeleteClient(clientID uint64) error {
	endpoint := fmt.Sprintf("/api/client/%d", clientID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// Создаем клиент
	client := NewWgEasyClient("http://localhost:51821", "Roman", "5N!xkD!Z4a!BRfv")

	// Получаем список клиентов
	clients, err := client.GetClients()
	if err != nil {
		log.Fatalf("Ошибка получения клиентов: %v", err)
	}
	fmt.Println("Список клиентов:")
	for _, c := range clients {
		fmt.Printf("- %s (ID: %d, Enabled: %v)\n", c.Name, c.ID, c.Enabled)
	}

	// Создаем нового клиента
	createClientName := fmt.Sprintf("new-client-%v", time.Now().Unix())
	_, err = client.CreateClient(createClientName)
	clientsNew, _ := client.GetClientsByName(createClientName)
	newClient := clientsNew[0]

	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}
	fmt.Printf("\nСоздан клиент: %s (ID: %d)\n", newClient.Name, newClient.ID)

	// Удаляем клиента
	if err := client.DeleteClient(newClient.ID); err != nil {
		log.Fatalf("Ошибка удаления клиента: %v", err)
	}
	fmt.Printf("\nКлиент %s удален\n", newClient.Name)
}
