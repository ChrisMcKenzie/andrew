package action

import (
	"context"

	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

type HandlerFunc func(context.Context, apiai.WebhookRequest) (*apiai.Fulfillment, error)

type Handler interface {
	Handle(context.Context, apiai.WebhookRequest) (*apiai.Fulfillment, error)
}

var actions map[string]Handler

type handler struct {
	f HandlerFunc
}

func (h *handler) Handle(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	return h.f(ctx, r)
}

func init() {
	actions = make(map[string]Handler)
}

func HandleFunc(action string, h HandlerFunc) {
	actions[action] = &handler{h}
}

func Handle(action string, h Handler) {
	actions[action] = h
}

func handleFunc(f HandlerFunc) Handler {
	return &handler{f}
}

func Request(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	handler, ok := actions[r.Result.Action]
	if !ok {
		handler = handleFunc(fallbackHandler)
	}

	return handler.Handle(ctx, r)
}

func fallbackHandler(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	return &r.Result.Fulfillment, nil
}
