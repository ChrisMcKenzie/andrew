package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ChrisMcKenzie/andrew/fulfillment/action"
	"github.com/ChrisMcKenzie/andrew/fulfillment/action/automatic"
	_ "github.com/ChrisMcKenzie/andrew/fulfillment/action/ifttt"
	_ "github.com/ChrisMcKenzie/andrew/fulfillment/action/math"
	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
	"github.com/boltdb/bolt"
)

var response = `{
	"speech": "idk",
	"displayText": "idk",
	"source": "andrew"
}`

var automaticAuthTemplate = "https://accounts.automatic.com/oauth/authorize/?client_id=%s&response_type=code&scope=scope:trip%%20scope:location%%20scope:vehicle:profile%%20scope:vehicle:events"

type oAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	Expires      int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("andrew.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	setupDB(db)

	http.Handle("/", http.FileServer(http.Dir("./ui")))

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

		f, err := action.Request(r.Context(), req)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", f)

		err = json.NewEncoder(w).Encode(f)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/oauth/automatic", oauthHandler)
	http.HandleFunc("/oauth/automatic/redirect", oauthAccessToken(db))

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Authorizations"))
		v := b.Get([]byte("automatic"))

		var results oAuthTokenResponse
		json.Unmarshal(v, &results)

		action.Handle("car.fuel-level", automatic.NewHandler(results.AccessToken, results.RefreshToken, db))
		return nil
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(
		w,
		r,
		fmt.Sprintf(automaticAuthTemplate, os.Getenv("AUTOMATIC_CLIENT_ID")),
		http.StatusTemporaryRedirect)
	return
}

func oauthAccessToken(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		resp, err := http.PostForm("https://accounts.automatic.com/oauth/access_token", url.Values{
			"client_id":     {os.Getenv("AUTOMATIC_CLIENT_ID")},
			"client_secret": {os.Getenv("AUTOMATIC_CLIENT_SECRET")},
			"code":          {code},
			"grant_type":    {"authorization_code"},
		})
		if err != nil {
			log.Printf("unable to authorize user: %s", err)
		}

		dec := json.NewDecoder(resp.Body)
		var response oAuthTokenResponse
		err = dec.Decode(&response)
		if err != nil {
			log.Printf("unable to authorize user: %s", err)
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Authorizations"))
			data, err := json.Marshal(response)
			if err != nil {
				return err
			}
			err = b.Put([]byte("automatic"), data)
			return err
		})

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func setupDB(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("Authorizations"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		return nil
	})
}
