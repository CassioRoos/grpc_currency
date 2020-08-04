package services

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// main struct
type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// Struct the will be get from XML
type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Cubes struct {
	// This brakes the XML
	// <Cubes first_obj>
	// 		<Cubes second_obj>
	// 			<Cubes our_obj CURRENCY, RATES>
	// 			</Cubes our_obj>
	// 		</Cubes second_obj>
	// </Cubes first_obj>

	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{l, map[string]float64{}}

	err := er.getRates()

	return er, err
}

// The way to get the currency for other currencies different from euro
func (e *ExchangeRates) GetRate(base, destination string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", base)
	}
	dr, ok := e.rates[destination]
	if !ok {
		return 0, fmt.Errorf("rate not found for currency %s", destination)
	}
	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200 got: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)
	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = r
	}
	e.rates["BRL"] = 1
	return nil
}

// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// just add a random difference to the rate and return it
				// this simulates the fluctuations in currency rates
				for k, v := range e.rates {
					// change can be 10% of original value
					change := (rand.Float64() / 10)
					// is this a postive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}

					// modify the rate
					e.rates[k] = v * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ret <- struct{}{}
			}
		}
	}()

	return ret
}
