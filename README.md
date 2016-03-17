# GoWebServer
Making a basic Go web server from [howistart] (https://howistart.org/posts/go/1).

Create a webserver to consume Open Weather Map and Weather Underground APIs to get average temps by cities on demand.

To use: 
Requires your own API Key for both APIS as environment variables. 

To use env vars in Golang requires you to import the OS package and then declare your vars in the function they will be used. *Note Golang doesn't like wasteful code, if you don't use it it will warn you in build.

Run on port 8080 or curl with /weather/ and city name to return temp in K.
