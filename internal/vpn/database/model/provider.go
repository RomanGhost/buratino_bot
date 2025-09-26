package model

type Provider struct {
	Name string `gorm:"size:16;primaryKey"`
}

var (
	Outline   = Provider{"outline"}
	Wireguard = Provider{"wireguard"}
)
