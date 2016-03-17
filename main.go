package main

//import everything we need
import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

func main() {
	//move our hello handler from root path to only /hello
	http.HandleFunc("/hello", hello)
	// listen on port 8080

	//we create our http handler function, 1st arg is a pattern bc it's a multiplexer that is going to match to known patterns
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		//how we get our 'city'. We split the url path on the slash, get 3 elements? pick the second one
		city := strings.SplitN(r.URL.Path, "/", 3)[2]

		data, err := query(city)
		//how we handle errors. if we get one...
		if err != nil {
			//we write to our client the error and then exit out
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//otherwise if successful we write to client w/header set to the following we give the client our json data
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

//similar to node
func hello(w http.ResponseWriter, r *http.Request) {
	//convert our string to bytes and when you receive an http request, write this to response
	w.Write([]byte("hello!"))
}

//populate using our struct and the API with a func
//we have a func query that takes in a city as a string and returns our weatherData structure AND error
func query(city string) (weatherData, error) {
	OpenWeatherApiKey := os.Getenv("OpenWeather_API_KEY")
	WeatherUndergroundApiKey := os.Getenv("Weather_Underground_API_KEY")

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + OpenWeatherApiKey + "&q=" + city)
	//if our request fails for some reason, i.e., it's not nil then return to us this err
	if err != nil {
		return weatherData{}, err
	}
	//if our GET is successful we don't close our resp body
	defer resp.Body.Close()
	//we then declare a var, d
	var d weatherData
	//call decode on our resp body with var d struct
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		//if something screws up? returns to us our new weatherData w/ error
		return weatherData{}, err
	}
	//otherwise if everything is fine, return our var d and error is nil
	return d, nil
}

//create a structure to get what we want from the API
//we define a type, declare it as a struct
type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
		//json is our tag allowing us to accss the encoding/json package
	} `json:"main"`
}

//let's create a general provider interface for multiple apis, OpenWeather/WeatherUnderground–whatever
type weatherProvider interface {
	temperature(city, string) (float64, error)
}

//create a new structure just for openWeatherMap data
new openWeatherMap struct{}
func (w openWeatherMap) temperature(city string) (float64, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID="+ OpenWeatherApiKey +"&q=" + city)
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
	if err:= json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	//log out the following string with city and the var d's main struct Kelvin
	log.Printf("openWeatherMap: %s: %.2f", city, d.Main.Kelvin)
	return d.Main.Kelvin, nil
}
//let's create a sp. type just for Weather Underground API
type WeatherUnderground struct {
	WeatherUndergroundApiKey string
}
//make a function to handle our Weather Underground API data
func (w WeatherUnderground) temperature(city string) (float64, error) {
	resp, err := http.GET("http://api.wunderground.com/api/" + w.WeatherUndergroundApiKey + "/conditions/q/" + city + ".json")
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
	if err:= json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return 0, err
	}
	//define kelvin which is var d's struct's Celsius and convert and return it
	kelvin := d.Observation.Celsius + 273.15
	log.Printf("weatherUnderground: %s: %.2f", city, kelvin)
	return kelvin, nil
}

//now that we have both APIs taken care of average out the temps and return it
func temperature(city string, providers ...weatherProvider) (float64, error) {
	//sum is our 'counter' var
	sum := 0.0

	//for loop--the only loop in Go. We memo out the initializer and put in provider as the range
	for _, provider := range providers {
		//take each provider's temp for a city param
		k, err := provider.temperature(city)
		//if we hit an error exit out return false/error
		if err != nil {
			return 0, err
		}
		//otherwise if ok, we add the kelvin temp for that city to sum
		sum += k
	}
	//when done with loop we avg out the temps and return it
	return sum / float64(len(providers)), nil
}

//create a func in main to query both apis and output formatted data
func main () {
	//not sure...assign a var with the a sp. interface and the two apis...?
	mw := multiWeatherProvider {
		openWeatherMap{},
		weatherUnderground{},
	}

	//write our handle func for the weather resource and our call back type func
	http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
		//define some vars we'll need
		begin := time.Now()
		city := strings.SplitN(r.url.Path, "/", 3)[2]
		//we get the temp for each city we request
		temp, err := mw.temperature(city)
		//our error handling
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//otherwise if ok we do the following: we set our headers
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//we Encode not Decode here, bc we're writing JSON mapping over every interface: our 2 apis in teh following format
		json.NewEncoder(w).Encode(map[string]interface{}{
			"city": city,
			"temp": temp, 
			"took": time.Since(begin).String(),
			})
		})
	//self explanatory
	http.ListenAndServe(":8080", nil)
}
