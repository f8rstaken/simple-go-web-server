package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

// Read the incoming data to the /trucks endpoint
// and return a Truck struct, where missing fields are nil
func readPostTruckData(r *http.Request) Truck {
	var plate *string
	var speed *int64
	var x *float64
	var y *float64

	if r.FormValue("plate") != "" {
		plateInput := r.FormValue("plate")
		plate = &plateInput
	}
	if r.FormValue("speed") != "" {
		speedInput, _ := strconv.ParseInt(r.FormValue("speed"), 10, 64)
		speed = &speedInput
	}
	if r.FormValue("x") != "" {
		xVal, _ := strconv.ParseFloat(r.FormValue("x"), 64)
		x = &xVal
	}
	if r.FormValue("y") != "" {
		yVal, _ := strconv.ParseFloat(r.FormValue("y"), 64)
		y = &yVal
	}

	return Truck{
		plate: plate,
		speed: speed,
		x:     x,
		y:     y,
	}

}

// Function used for handling the /trucks endpoint
func trucksHandler(dbInstance *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var truck Truck
		if r.Method != "GET" {
			if err := r.ParseForm(); err != nil {
				fmt.Errorf("Error while parsing data to /trucks endpoint: %v", err)
				return
			}
			truck = readPostTruckData(r)
		}

		switch r.Method {
		case "GET":
			var trucks, err = queryTrucks(dbInstance)
			if err != nil {
				fmt.Printf("GET /trucks: %v", err)
				return
			}
			fmt.Fprintf(w, "GET Trucks: %v", trucks)
		case "POST":
			_, err := insertTruck(dbInstance, truck)
			if err != nil {
				fmt.Printf("POST /trucks: %v", err)
				return
			}
		case "PATCH":
			err := updateTruck(dbInstance, truck)
			if err != nil {
				fmt.Printf("PATCH /trucks: %v", err)
				return
			}
		case "DELETE":
			err := deleteTruck(dbInstance, *truck.plate)
			if err != nil {
				fmt.Printf("DELETE /trucks: %v", err)
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
