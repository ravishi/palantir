package bennu

type envelope struct {
	Ref string `json:"ref"`
	Topic string `json:"topic"`
	Event string `json:"event"`
	Payload interface{} `json:"payload"`
}
