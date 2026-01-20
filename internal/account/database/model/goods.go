package model

import "gorm.io/gorm"

type GoodsPrice struct {
	gorm.Model
	SysName string `gorm:"type:varchar(128);unique;not null"`
	Name    string `gorm:"type:varchar(128);unique;not null"`
	Price   int64  `gorm:"not null"`
}

var (
	VPN1Min   = GoodsPrice{SysName: "1m vpn", Name: "1 минута", Price: 5}         // 216 rub/month	0.005 rub
	VPN1Hour  = GoodsPrice{SysName: "1h vpn", Name: "1 час", Price: 250}          // 180 rub/month	0.25 rub
	VPN1Day   = GoodsPrice{SysName: "1d vpn", Name: "1 день", Price: 5000}        // 150 rub/month	5 rub
	VPN1Month = GoodsPrice{SysName: "1month vpn", Name: "1 месяц", Price: 140000} // 140 rub/month	140 rub
	TopUP     = GoodsPrice{SysName: "topUP", Name: "Пополнение", Price: -1}       // 0.001(10^-3) rub
)
