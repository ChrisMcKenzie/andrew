package monitor

type Notification struct {
	Message  string `json:"message"`
	Provider string `json:"provider"`
}

type Monitor interface {
	Monitor(chan<- Notification)
}
