package automatic

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChrisMcKenzie/andrew/fulfillment/apiai"
	"github.com/boltdb/bolt"
)

type Vehicle struct {
	URL              string    `json:"url"`
	ID               string    `json:"id"`
	VIN              string    `json:"vin"`
	Created          time.Time `json:"created_at"`
	Updated          time.Time `json:"updated_at"`
	Make             string    `json:"make"`
	Model            string    `json:"model"`
	Year             int       `json:"year"`
	SubModel         string    `json:"submodel"`
	DisplayName      string    `json:"display_name"`
	FuelGrade        string    `json:"fuel_grade"`
	FuelLevelPercent float64   `json:"fuel_level_percent"`
	BatteryVoltage   float64   `json:"battery_voltage"`
	ActiveDTCs       []struct {
		Code        string    `json:"code"`
		Created     time.Time `json:"created_at"`
		Description string    `json:"description"`
	} `json:"active_dtcs"`
}

type AutomaticHandler struct {
	*http.Client
	db                        *bolt.DB
	AccessToken, RefreshToken string
}

func NewHandler(access, refresh string, db *bolt.DB) *AutomaticHandler {
	ah := &AutomaticHandler{&http.Client{}, db, access, refresh}

	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Vehicles"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		cars, err := ah.getVehicles()
		if err != nil {
			return err
		}

		for _, car := range cars {
			data, err := json.Marshal(car)
			if err != nil {
				return err
			}
			err = b.Put([]byte(car.Model), data)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("unable to create vehicle list: %s", err)
	}

	return ah
}

func (a *AutomaticHandler) Handle(ctx context.Context, r apiai.WebhookRequest) (*apiai.Fulfillment, error) {
	carName := r.Result.Parameters["car"].(string)

	var car Vehicle
	err := a.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Vehicles"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}
		val := b.Get([]byte(carName))

		if val != nil {
			err := json.Unmarshal(val, &car)
			if err != nil {
				return err
			}
		}
		fmt.Printf("%+v", car)

		return nil
	})

	if err != nil {
		return &apiai.Fulfillment{
			Speech: fmt.Sprintf("I was unable to get the fuel level for your  %s because %s", carName, err),
		}, err
	}

	v, err := a.getVehicle(car.URL)
	if err != nil {
		return &apiai.Fulfillment{
			Speech: fmt.Sprintf("I was unable to get the fuel level for your  %s because %s", carName, err),
		}, err
	}

	return &apiai.Fulfillment{
		Speech: fmt.Sprintf("Your %s's fuel level is at %.1f%%", v.Model, v.FuelLevelPercent),
	}, nil
}

func (a *AutomaticHandler) getVehicle(url string) (*Vehicle, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken))

	resp, err := a.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		dec := json.NewDecoder(resp.Body)

		var result Vehicle
		err := dec.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	}

	return nil, fmt.Errorf("unable to retrieve vehicle: %s", resp.Status)
}

func (a *AutomaticHandler) getVehicles() ([]Vehicle, error) {
	req, err := http.NewRequest("GET", "https://api.automatic.com/vehicle/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.AccessToken))

	resp, err := a.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		dec := json.NewDecoder(resp.Body)

		var results struct {
			Results []Vehicle `json:"results"`
		}
		err := dec.Decode(&results)
		if err != nil {
			return nil, err
		}

		return results.Results, nil
	}

	return nil, fmt.Errorf("unable to retrieve vehicles")
}
