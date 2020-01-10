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
	ID         int    `json:"id"`
	HardwareId string `json:"hardwareId"`
}

type HardwareId struct {
	HardwareId string `json:"hardwareId"`
}

type SoftId struct {
	SoftId int `json:"softId"`
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
	http.HandleFunc("/sendMeasurements", sendMeasurements)
	dbConn()

	log.Println("Listening on port 4040")
	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatalln(err)
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

	var hardwareId HardwareId
	//var hardwareID string
	var deviceId int

	err := json.NewDecoder(r.Body).Decode(&hardwareId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Fprintln(w, "JSON error")
		return
	}

	deviceId = findDeviceInDb(hardwareId.HardwareId)

	if deviceId == 0 {
		addDeviceToDb(hardwareId.HardwareId)
		deviceId = findDeviceInDb(hardwareId.HardwareId)
	}

	var softId SoftId
	softId.SoftId = deviceId

	b, _ := json.Marshal(softId) //Convert slice to JSON byte array
	s := string(b)               //JSON byte array to string
	fmt.Fprintln(w, s)           //Print to browser

}

func findDeviceInDb(hardwareId string) int {
	db := dbConn()
	getDeviceId, err := db.Query("Select id from Device where HardwareId = ?", hardwareId)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	defer getDeviceId.Close()

	var foundDeviceId int
	for getDeviceId.Next() {
		err := getDeviceId.Scan(&foundDeviceId)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Found id: ", foundDeviceId)
	}

	return foundDeviceId
}

func addDeviceToDb(hardwareId string) {
	fmt.Println("Adding device to database, hardwareId: ", hardwareId)
	db := dbConn()
	addDeviceToDb, err := db.Query("INSERT INTO Device VALUES(?, ?)", 0, hardwareId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer addDeviceToDb.Close()
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
	defer insertMeasurement.Close()
	insertMeasurement.Exec(0, measurement.ID, measurement.CO2Value, measurement.TVOCValue, measurement.TimeStamp)
}
