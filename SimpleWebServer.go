package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var randomWords = []string{"Raspberry", "ESP32", "Lopy4", "Arduino"}
	var ids = []int{0}

	http.HandleFunc("/4", handler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		fmt.Fprintln(w, "Hello, world!")
	})

	http.HandleFunc("/getNewId", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		id := ids[len(ids)-1]
		fmt.Fprintln(w, id)
		ids = append(ids, id+1)
	})

	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		b, _ := json.Marshal(randomWords)
		s := string(b)
		fmt.Fprintln(w, s)
	})

	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["key"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	log.Println("Url Param 'key' is: " + string(key))
}
