package ifttt

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ChrisMcKenzie/andrew/fulfillment/action"
	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

var makerTemplate = fmt.Sprintf("http://maker.ifttt.com/trigger/%%s/with/key/%s", os.Getenv("MAKER_TOKEN"))

func init() {
	action.HandleFunc("lights.on", lightsOnHandler)
	action.HandleFunc("lights.off", lightsOffHandler)
}

func lightsOnHandler(ctx context.Context, req apiai.WebhookRequest) (*apiai.Fulfillment, error) {

	for _, location := range req.Result.Parameters["location"].([]interface{}) {
		_, err := http.Get(fmt.Sprintf(makerTemplate,
			fmt.Sprintf("%s_lights_on", location)))
		if err != nil {
			return nil, err
		}
	}

	return &req.Result.Fulfillment, nil
}

func lightsOffHandler(ctx context.Context, req apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	for _, location := range req.Result.Parameters["location"].([]interface{}) {
		_, err := http.Get(fmt.Sprintf(makerTemplate,
			fmt.Sprintf("%s_lights_off", location)))
		if err != nil {
			return nil, err
		}
	}

	return &req.Result.Fulfillment, nil
}
