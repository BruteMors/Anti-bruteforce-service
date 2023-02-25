package entity

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
}
