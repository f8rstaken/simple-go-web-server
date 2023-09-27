package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const DB_USERNAME = "root"
const DB_PASSWORD = "root"
const DB_ADDR = "127.0.0.1:3306"

type Truck struct {
	plate *string
	speed *int64
	x     *float64
	y     *float64
}

type TruckData struct {
	Plate string  `json:"plate"`
	Speed int64   `json:"speed"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

func connectToDatabase(username string, password string, addr string) *sql.DB {

	cfg := mysql.Config{
		User:   username,
		Passwd: password,
		Addr:   addr,
		Net:    "tcp",
		DBName: "dmt",
	}

	var db *sql.DB
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}

func queryTrucks(dbInstance *sql.DB) ([]TruckData, error) {
	var trucks []TruckData

	rows, err := dbInstance.Query("SELECT * FROM trucks")
	if err != nil {
		return nil, fmt.Errorf("GET /trucks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var truck TruckData
		if err := rows.Scan(&truck.Plate, &truck.Speed, &truck.X, &truck.Y); err != nil {
			return nil, fmt.Errorf("GET /trucks: %v", err)
		}
		trucks = append(trucks, truck)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GET /trucks: %v", err)
	}
	return trucks, nil
}

func convertToJSON(trucks []TruckData) ([]byte, error) {
	jsonData, err := json.Marshal(trucks)
	if err != nil {
		return nil, fmt.Errorf("Converting to JSON: %v", err)
	}
	return jsonData, nil
}

func insertTruck(dbInstance *sql.DB, truck Truck) (int64, error) {
	result, err := dbInstance.Exec("INSERT INTO trucks (plate, speed, x, y) VALUES (?, ?, ?, ?)", *truck.plate, *truck.speed, *truck.x, *truck.y)
	if err != nil {
		return -1, fmt.Errorf("POST /trucks: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("POST /trucks: %v", err)
	}
	return id, nil
}

func updateTruck(dbInstance *sql.DB, truck Truck) error {
	if truck.plate == nil || *truck.plate == "" {
		return fmt.Errorf("Provided truck with no plate")
	}

	updatedFields := []string{}
	var values []interface{}

	updatedFields = append(updatedFields, "plate = ?")
	values = append(values, *&truck.plate)

	if truck.speed != nil {
		updatedFields = append(updatedFields, "speed = ?")
		values = append(values, *&truck.speed)
	}

	if truck.x != nil {
		updatedFields = append(updatedFields, "x = ?")
		values = append(values, *&truck.x)
	}

	if truck.y != nil {
		updatedFields = append(updatedFields, "y = ?")
		values = append(values, *&truck.y)
	}
	values = append(values, *&truck.plate)

	updateQuery := "UPDATE trucks SET " + strings.Join(updatedFields, ", ") + " WHERE plate = ?"
	_, err := dbInstance.Exec(updateQuery, values...)
	if err != nil {
		return fmt.Errorf("PATCH /trucks: %v", err)
	}
	return nil
}

func deleteTruck(dbInstance *sql.DB, plate string) error {
	_, err := dbInstance.Exec("DELETE FROM Trucks WHERE plate = ?", plate)
	if err != nil {
		return fmt.Errorf("DELETE /trucks: %v", err)
	}
	return nil
}
