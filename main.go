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
