package math

import (
	"context"
	"fmt"

	"github.com/ChrisMcKenzie/andrew/fulfillment/action"
	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

func init() {
	action.HandleFunc("math.simple", doSimpleMath)
}

func doSimpleMath(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {

	var answer float64
	switch r.Result.Parameters["action"].(string) {
	case "plus":
		answer = r.Result.Parameters["number"].(float64) + r.Result.Parameters["number1"].(float64)
	}

	return &apiai.Fulfillment{
		Speech: fmt.Sprintf("the answer %d", answer),
	}, nil
}
