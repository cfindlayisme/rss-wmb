package model

type WebhookMessage struct {
	Message    string
	Password   string
	ColourCode *int8 // https://modern.ircdocs.horse/formatting.html
}

type DirectedWebhookMessage struct {
	IncomingMessage WebhookMessage
	Target          string
}
