package types

type WSMessage struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Login struct {
	ClientID int    `json:"clientID"`
	UserName string `json:"username"`
}
