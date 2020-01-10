package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Device struct {
	ID       int `json:"id"`
	DeviceID int `json:"deviceid"`
}

type MeasurementCollection struct {
	ID   int   `json:"id"`
	ECO2 []int `json:"eCO2"`
	TVOC []int `json:"TVOC"`
}

type Measurement struct {
	ID        int       `json:"id"`
	CO2Value  int       `json:"co2values"`
	TVOCValue int       `json:"tvocvalues"`
	TimeStamp time.Time `json:"timestamp"`
}

var (
	ctx context.Context
	db  *sql.DB
)

func main() {
	http.HandleFunc("/getNewId", getNewId)
	http.HandleFunc("/respondWithJson", respondWithJSON)
	http.HandleFunc("/OpenDB", testDatabase)
	http.HandleFunc("/sendMeasurements", sendMeasurements)
	dbConn()

	log.Println("Listening on port 4040")
	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatalln(err)
	}
}

func testDatabase(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
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

func sendMeasurements(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	// Declare a new Measurement struct
	var measurementCollection MeasurementCollection

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 statys code.
	err := json.NewDecoder(r.Body).Decode(&measurementCollection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Code for handling the Measurement
	for i := 0; i < len(measurementCollection.ECO2); i++ {
		var measurement Measurement
		measurement.ID = measurementCollection.ID
		measurement.CO2Value = measurementCollection.ECO2[i]
		measurement.TVOCValue = measurementCollection.TVOC[i]
		measurement.TimeStamp = time.Now()

		addMeasurementToDatabase(measurement)
	}
}

func addMeasurementToDatabase(measurement Measurement) {
	fmt.Println("Adding measurement to db, id value: ", measurement.ID, " eCO2 value: ", measurement.CO2Value, " TVOC value: ", measurement.TVOCValue, " Date value: ", measurement.TimeStamp)
	db := dbConn()
	insertMeasurement, err := db.Prepare("INSERT INTO Readings VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
	}
	insertMeasurement.Exec(0, measurement.ID, measurement.CO2Value, measurement.TVOCValue, measurement.TimeStamp)
}
