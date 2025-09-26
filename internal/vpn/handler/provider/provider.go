package provider

import "github.com/RomanGhost/buratino_bot.git/internal/vpn/handler/provider/data"

type Provider interface {
	CreateKey(name string) (*data.KeyConnectData, error)
	// GetKeyByName(name string) (*data.AccessKey, error)
	DeleteAccessKey(keyID int) error
}
