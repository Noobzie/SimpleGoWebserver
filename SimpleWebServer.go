package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Device struct {
	ID       int `json:"id"`
	DeviceID int `json:"deviceid"`
}

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	http.HandleFunc("/4", readFromUrl) //Call function outside of main
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/getNewId", getNewId)
	http.HandleFunc("/3", respondWithJSON)
	http.HandleFunc("/OpenDB", testDatabase)
	dbConn()

	log.Println("Listening on port 4040")
	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatalln(err)
	}
}

func testDatabase(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	results, err := db.Query("Select id, deviceid from Device")

	for results.Next() {
		var device Device
		err = results.Scan(&device.ID, &device.DeviceID)
		if err != nil {
			fmt.Fprintln(w, "Some kind of error")
		}
		fmt.Fprintln(w, device.ID)
		fmt.Fprintln(w, device.DeviceID)
	}

	if err != nil {
		fmt.Fprintln(w, "Some kind of error")
		fmt.Fprintln(w, err)
	}

}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "dasc"
	dbPass := "dasc"
	dbName := "CO2-database"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func helloWorld(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080
	log.Println(r.Method, r.URL.Path)
	fmt.Fprintln(w, "Hello, world!") //Print to browser
}

func getNewId(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/getNewId, gives the last value in the slice and adds last value + 1 to the slice
	log.Println(r.Method, r.URL.Path)
	deviceId := 1345
	id := 999999
	db := dbConn()
	results, err := db.Query("Select id, deviceid from Device")

	for results.Next() {
		var device Device
		err = results.Scan(&device.ID, &device.DeviceID)
		if err != nil {
			fmt.Fprintln(w, "Some kind of error")
		}
		if device.DeviceID == deviceId {
			id = device.ID
		}
	}

	if id == 999999 {
		insertDevice, err := db.Prepare("INSERT INTO Device(DeviceID) VALUES(?)")
		if err != nil {
			panic(err.Error())
		}
		insertDevice.Exec(deviceId)

		results, err := db.Query("Select id, deviceid from Device")

		for results.Next() {
			var device Device
			err = results.Scan(&device.ID, &device.DeviceID)
			if err != nil {
				fmt.Fprintln(w, "Some kind of error")
			}
			if device.DeviceID == deviceId {
				id = device.ID
				fmt.Fprintln(w, "DeviceId: ", id)
			}
		}
	} else {
		fmt.Fprintln(w, "DeviceId: ", id)
	}
}

func respondWithJSON(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/3 with a JSON example
	var randomWords = []string{"Raspberry", "ESP32", "Lopy4", "Arduino"} //Slice of strings
	log.Println(r.Method, r.URL.Path)
	b, _ := json.Marshal(randomWords) //Convert slice to JSON byte array
	s := string(b)                    //JSON byte array to string
	fmt.Fprintln(w, s)                //Print to browser
}

func readFromUrl(w http.ResponseWriter, r *http.Request) { //responds to localhost:8080/4 and finds the parameters in the url.		Example: localhost:8080/?key=CO2%20655	<- finds CO2 and 655

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
