package provider

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/RomanGhost/buratino_bot.git/internal/vpn/database/model"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider/data"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider/outline"
	"github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider/wireguard"
)

type Provider interface {
	CreateKey(name string) (*data.KeyConnectData, error)
	// GetKeyByName(name string) (*data.AccessKey, error)
	DeleteAccessKey(keyID int) error
}

func NewProvider(link string, providerName string) Provider {
	switch providerName {
	case model.Outline.Name:
		return createOutlineProvider(link)
	case model.Wireguard.Name:
		return createWireguardProvider(link)
	default:
		return nil
	}
}

func createOutlineProvider(link string) Provider {
	return outline.NewOutlineClient(link)
}

func createWireguardProvider(link string) Provider {
	// link strcut https://dns:port/base64(login:password)

	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	baseURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	// Берем строку после "/"
	encoded := strings.TrimPrefix(u.Path, "/")

	// Декодируем base64
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(err)
	}
	decoded := string(decodedBytes)

	// Делим на логин и пароль
	var login, password string
	if parts := strings.SplitN(decoded, ":", 2); len(parts) == 2 {
		login = parts[0]
		password = parts[1]
	}
	return wireguard.NewWgEasyClient(baseURL, login, password)
}
