package main

//import everything we need
import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//create a func in main to query both apis and output json formatted data
func main() {
	//not sure about this...assigning  to 'mw' our structs in func?
	mw := multiWeatherProvider{
		openWeatherMap{},
		weatherUnderground{},
	}

	//write our handle func for the weather resource and our call back type func...asterisk means?
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		//define some vars we'll need
		begin := time.Now()
		city := strings.SplitN(r.URL.Path, "/", 3)[2]
		//we get the temp for each city we request
		temp, err := mw.temperature(city)
		//our error handling
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//otherwise if ok we do the following: we set our headers
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//we Encode not Decode here, bc we're writing JSON mapping over every interface: our 2 apis in the following format
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city": city,
			"temp": temp,
			"took": time.Since(begin).String(),
		})
	})
	//self explanatory
	http.ListenAndServe(":8080", nil)
}

//declare our interface that we need to satisfy with the data from an API
type weatherProvider interface {
	temperature(city string) (float64, error)
}

//?? declare a type for all
type multiWeatherProvider []weatherProvider

//now that we have both APIs taken care of average out the temps and return it, using concurrency
func (w multiWeatherProvider) temperature(city string) (float64, error) {
	//we make two channels one for temps one for errors with a type and length
	temps := make(chan float64, len(w))
	errs := make(chan error, len(w))

	//for loop--the only loop in Go. We memo out the initializer and put in provider as the range
	for _, provider := range w {
		go func(p weatherProvider) {
			k, err := p.temperature(city)
			if err != nil {
				errs <- err
				return
			}
			temps <- k
		}(provider)
	}
	//a counter var
	sum := 0.0
	//loop through each json object
	for i := 0; i < len(w); i++ {
		//Rob Pike says that without select keyword its not concurrency
		select {
		case temp := <-temps:
			sum += temp
		case err := <-errs:
			return 0, err
		}
	}
	return sum / float64(len(w)), nil
}

//create a new structure just for openWeatherMap data
type openWeatherMap struct{}

func (w openWeatherMap) temperature(city string) (float64, error) {
	OpenWeatherApiKey := os.Getenv("OpenWeather_API_KEY")

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + OpenWeatherApiKey + "&q=" + city)
	//if we encounter an error exit out with the error
	if err != nil {
		return 0, err
	}
	//if we are 200 OK don't close
	defer resp.Body.Close()
	//define a var d structure that'll hold our resp data
	var d struct {
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	//now decode that json resp body and if we encounter an error exit out w/error
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	//log out the following string with city and the var d's main struct Kelvin
	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}

//let's create a sp. type just for Weather Underground API
type weatherUnderground struct {
	//why is this missing from Open Weather?
	WeatherUndergroundApiKey string
}

//make a function to handle our Weather Underground API data
func (w weatherUnderground) temperature(city string) (float64, error) {
	WeatherUndergroundApiKey := os.Getenv("Weather_Underground_API_KEY")

	resp, err := http.Get("http://api.wunderground.com/api/" + WeatherUndergroundApiKey + "/conditions/q/" + city + ".json")
	//if we hit an error exit out with false and return to me the error
	if err != nil {
		return 0, err
	}

	//otherwise if 200 ok don't close the response body
	defer resp.Body.Close()

	//create a var d structure to hold our data, going to convert that Celsius to K
	var d struct {
		Observation struct {
			Celsius float64 `json:"temp_c"`
		} `json:"current_observation"`
	}
	//Decode and if you hit an error return it
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	//define kelvin which is var d's struct's Celsius and convert and return it
	kelvin := d.Observation.Celsius + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}
