package types

type Client struct {
	Id         string `json:"id" gorm:"primaryKey"`
	ClientType int    `json:"clienttype"` //接入的客户端类型，
	Token      string `json:"token"`      //token
	AccessID   string `json:"accessid"`   //接入的客户端ID
	Addr       string `json:"addr"`       //接入的客户端地址
	Os         string `json:"os"`         //接入的客户端操作系统
	Info       string `json:"info"`       //一些接入信息的描述，比如硬件信息
}
