package main

//import everything we need
import (
	"encoding/json"
	"net/http"
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
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YourAPIKeyHere=" + city)
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
