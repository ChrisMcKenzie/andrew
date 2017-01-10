package math

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ChrisMcKenzie/andrew/fulfillment/action"
	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

func init() {
	action.HandleFunc("math.simple", doSimpleMath)
}

func doSimpleMath(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	var answer float64
	number, err := strconv.ParseFloat(r.Result.Parameters["number"].(string), 64)
	number1, err := strconv.ParseFloat(r.Result.Parameters["number1"].(string), 64)
	if err != nil {
		return nil, err
	}

	switch r.Result.Parameters["math-action"].(string) {
	case "plus":
		answer = number + number1
	case "minus":
		answer = number - number1
	case "multiply":
		answer = number * number1
	case "divide":
		answer = number / number1
	}

	return &apiai.Fulfillment{
		Speech: fmt.Sprintf("the answer %v", answer),
	}, nil
}
