package main

// we import the net/http package from the standard library
import "net/http"

func main() {
	// we create a function to handle at root path of our server and pass in our hello func below
	http.HandleFunc("/", hello)
	// listen on port 8080
	http.ListenAndServe(":8080", nil)
}

//similar to node
func hello(w http.ResponseWriter, r *http.Request) {
	//convert our string to bytes and when you receive an http request, write this to response
	w.Write([]byte("hello!"))
}

//create a structure to get what we want from the API
//we define a type, declare it as a struct
type weatherData struct {
    Main struct {
        Kelvin float64 'json:"temp"'
        //json is our tag allowing us to accss the encoding/json package
    } 'json:"main"'
}

//populate using our struct and the API with a func
//we have a func query that takes in a city as a string and returns our weatherData structure AND error
func query(city, string) (weatherData, error) {
    resp, err:= http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=YOUR_API_KEY&q=" + city)
    //if our request fails for some reason, i.e., it's not nil then return to us this err
    if err != nil {
        return weatherData{}, err
    }
    //if our GET is successful we don't close our resp body
    defer resp.Body.Close()
    //we then declare a var, d
    var d weatherData
    //call decode on our resp body with var d struct
    if err:= json.NewDecoder(resp.Body).Decode(&d);
        //returns to us our new weatherData w/ nil error
        return weatherData{}, err
}

return d, nil
