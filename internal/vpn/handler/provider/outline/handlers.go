package outline

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Структуры для API ответов
type AccessKey struct {
	ID        int    `json:"id,string"`
	Name      string `json:"name,omitempty"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessURL string `json:"accessUrl"`
	DataLimit *struct {
		Bytes int64 `json:"bytes"`
	} `json:"dataLimit,omitempty"`
	UsedBytes int64 `json:"usedBytes,omitempty"`
}

type AccessKeysResponse struct {
	AccessKeys []AccessKey `json:"accessKeys"`
}

type ServerInfo struct {
	Name               string `json:"name"`
	ServerId           string `json:"serverId"`
	MetricsEnabled     bool   `json:"metricsEnabled"`
	CreatedTimestampMs int64  `json:"createdTimestampMs"`
	Version            string `json:"version"`
	AccessKeyDataLimit *struct {
		Bytes int64 `json:"bytes"`
	} `json:"accessKeyDataLimit,omitempty"`
	PortForNewAccessKeys  int    `json:"portForNewAccessKeys"`
	HostnameForAccessKeys string `json:"hostnameForAccessKeys"`
}

// Клиент для работы с Outline API
type OutlineClient struct {
	baseURL    string
	httpClient *http.Client
}

// Создание нового клиента
func NewOutlineClient(apiURL string) *OutlineClient {
	// Создаем HTTP клиент с отключенной проверкой SSL
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &OutlineClient{
		baseURL:    apiURL,
		httpClient: client,
	}
}

// Выполнение HTTP запроса
func (c *OutlineClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
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

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}

	return resp, nil
}

// Получение списка всех ключей
func (c *OutlineClient) GetAccessKeys() ([]AccessKey, error) {
	resp, err := c.makeRequest("GET", "/access-keys", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	// bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Printf("Body: %+v\n", string(bodyBytes))
	var keysResponse AccessKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&keysResponse); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return keysResponse.AccessKeys, nil
}

// Создание нового ключа
func (c *OutlineClient) CreateAccessKey() (*AccessKey, error) {
	resp, err := c.makeRequest("POST", "/access-keys", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	var key AccessKey
	if err := json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &key, nil
}

// Удаление ключа
func (c *OutlineClient) DeleteAccessKey(keyID int) error {
	endpoint := fmt.Sprintf("/access-keys/%d", keyID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

// Переименование ключа
func (c *OutlineClient) RenameAccessKey(keyID int, name string) error {
	endpoint := fmt.Sprintf("/access-keys/%d/name", keyID)
	body := map[string]string{"name": name}

	resp, err := c.makeRequest("PUT", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

// Установка лимита трафика для ключа
func (c *OutlineClient) SetDataLimit(keyID int, limitBytes int64) error {
	endpoint := fmt.Sprintf("/access-keys/%d/data-limit", keyID)
	body := map[string]interface{}{
		"limit": map[string]int64{
			"bytes": limitBytes,
		},
	}

	resp, err := c.makeRequest("PUT", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

// Удаление лимита трафика для ключа
func (c *OutlineClient) RemoveDataLimit(keyID int) error {
	endpoint := fmt.Sprintf("/access-keys/%d/data-limit", keyID)

	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	return nil
}

// Получение информации о сервере
func (c *OutlineClient) GetServerInfo() (*ServerInfo, error) {
	resp, err := c.makeRequest("GET", "/server", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	var serverInfo ServerInfo
	if err := json.NewDecoder(resp.Body).Decode(&serverInfo); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return &serverInfo, nil
}

// Получение метрик использования данных
func (c *OutlineClient) GetDataUsage() (map[string]int64, error) {
	resp, err := c.makeRequest("GET", "/metrics/transfer", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	var usage map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return usage, nil
}

// Утилитарная функция для конвертации GB в байты
func GBToBytes(gb int) int64 {
	return int64(gb) * 1024 * 1024 * 1024
}

// Утилитарная функция для конвертации байтов в GB
func BytesToGB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
