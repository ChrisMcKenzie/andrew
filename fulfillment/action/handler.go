package action

import (
	"context"

	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

type HandlerFunc func(context.Context, apiai.WebhookRequest) (*apiai.Fulfillment, error)

var actions map[string]HandlerFunc

func init() {
	actions = make(map[string]HandlerFunc)
}

func Register(action string, h HandlerFunc) {
	actions[action] = h
}

func Handler(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	handler, ok := actions[r.Result.Action]
	if !ok {
		handler = fallbackHandler
	}

	return handler(ctx, r)
}

func fallbackHandler(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	return &r.Result.Fulfillment, nil
}
