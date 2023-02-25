package entity

type IpNetwork struct {
	Ip   string `json:"ip" db:"prefix"`
	Mask string `json:"mask" db:"mask"`
}
