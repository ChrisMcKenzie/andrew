package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ChrisMcKenzie/andrew/fulfillment/action"
	_ "github.com/ChrisMcKenzie/andrew/fulfillment/action/ifttt"
	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
)

var response = `{
	"speech": "idk",
	"displayText": "idk",
	"source": "andrew"
}`

func main() {
	http.HandleFunc("/intent", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)

		var req apiai.WebhookRequest
		err := dec.Decode(&req)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%+v\n", req)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		f, err := action.Handler(r.Context(), req)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = json.NewEncoder(w).Encode(f)
		if err != nil {
			fmt.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
