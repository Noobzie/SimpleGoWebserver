package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var randomWords = []string{"Raspberry", "ESP32", "Lopy4", "Arduino"} //Slice of strings
	var ids = []int{0}                                                   //Slice of int's

	http.HandleFunc("/4", handler) //Call function outside of main

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080
		log.Println(r.Method, r.URL.Path)
		fmt.Fprintln(w, "Hello, world!") //Print to browser
	})

	http.HandleFunc("/getNewId", func(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/getNewId, gives the last value in the slice and adds last value + 1 to the slice
		log.Println(r.Method, r.URL.Path)
		id := ids[len(ids)-1]   //Find the last int in slice
		fmt.Fprintln(w, id)     //Print to browser
		ids = append(ids, id+1) //Add new int on the end of the slice
	})

	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/3 with a JSON example
		log.Println(r.Method, r.URL.Path)
		b, _ := json.Marshal(randomWords) //Convert slice to JSON byte array
		s := string(b)                    //JSON byte array to string
		fmt.Fprintln(w, s)                //Print to browser
	})

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/4 and finds the parameters in the url.		Example: localhost:8080/?key=CO2%20655	<- finds CO2 and 655

	keys, ok := r.URL.Query()["key"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	log.Println("Url Param 'key' is: " + string(key)) //Print to console
}
