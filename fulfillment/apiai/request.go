package apiai

type Message struct {
	Speech   string      `json:"speech,omitempty"`
	Type     int         `json:"type"`
	Payload  interface{} `json:"payload,omitempty"`
	Title    string      `json:"title,omitempty"`
	SubTitle string      `json:"subtitle,omitempty"`
	Replies  []string    `json:"replies,omitempty"`
	ImageURL string      `json:"imageUrl,omitempty"`
	Buttons  []struct {
		Text     string `json:"text"`
		PostBack string `json:"postback"`
	} `json:"buttons,omitempty"`
}

type Fulfillment struct {
	Speech   interface{} `json:"speech"`
	Messages []Message   `json:"messages,omitempty"`
}

type Context struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
	Lifespan   int                    `json:"lifespan"`
}

type Result struct {
	Speech        string                 `json:"speech"`
	Fulfillment   Fulfillment            `json:"fulfillment"`
	Action        string                 `json:"action"`
	ResolvedQuery string                 `json:"resolvedQuery"`
	Parameters    map[string]interface{} `json:"parameters"`
	Contexts      []Context              `json:"contexts"`
}

type WebhookRequest struct {
	Result Result `json:"result"`
}
