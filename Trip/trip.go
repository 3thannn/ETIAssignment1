package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Trip struct {
	TripID       int
	PassengerID  int
	DriverID     int
	Pickup       int
	Dropoff      int
	Pickuptime   string
	Dropofftime  string
	Carlicenseno string
	Status       string
}

type Driver struct {
	DriverID     int
	Firstname    string
	Lastname     string
	Mobileno     string
	Emailaddress string
	Nric         string
	Carlicenseno string
	Availability int
}

type Passenger struct {
	PassengerID  int
	Firstname    string
	Lastname     string
	Mobileno     string
	Emailaddress string
}

//db function to get the current trip that the driver is in
func getdrivertrip(db *sql.DB, driverID string) Trip {
	query := fmt.Sprintf("SELECT TripID, DriverID, PassengerID, Pickup, Dropoff, ifnull(PickUpTime, '-'), ifnull(DropOffTime, '-'), CarLicenseNo, Status FROM Trip WHERE DriverID = %s AND Status != 'Completed'", driverID)
	results, err := db.Query(query)
	var trip Trip
	if err != nil {
		panic(err.Error())
	} else {
		for results.Next() { //go through all results to append to array
			// map this type to the record in the table
			results.Scan(&trip.TripID, &trip.DriverID, &trip.PassengerID, &trip.Pickup, &trip.Dropoff, &trip.Pickuptime, &trip.Dropofftime, &trip.Carlicenseno, &trip.Status)
		}
	}
	return trip
}

//http function to handle get requests for the current trip that the driver is in
func getDriverCurrentTrip(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	} else {
		params := mux.Vars(r)
		driverID := params["driverid"]
		defer db.Close()
		if r.Method == "GET" {
			if err == nil {
				var trip Trip = getdrivertrip(db, driverID)
				json.NewEncoder(w).Encode(trip)
			} else {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			}
		}
	}
}

//db function to get all passenger trips
func gettrips(db *sql.DB, passengerID string) []Trip {
	query := fmt.Sprintf("SELECT TripID, DriverID, PassengerID, Pickup, Dropoff, ifnull(PickUpTime, '-'), ifnull(DropOffTime, '-'), CarLicenseNo, Status FROM Trip WHERE PassengerID = %s ORDER BY PickUpTime DESC", passengerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var trips []Trip
	for results.Next() { //go through all results to append to array
		var trip Trip
		// map this type to the record in the table
		results.Scan(&trip.TripID, &trip.DriverID, &trip.PassengerID, &trip.Pickup, &trip.Dropoff, &trip.Pickuptime, &trip.Dropofftime, &trip.Carlicenseno, &trip.Status)
		trips = append(trips, trip)
	}
	return trips
}

//http function to handle get requests for retrieving passenger trips
func getpassengertrips(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	} else {
		params := mux.Vars(r)
		passengerID := params["passengerid"]
		defer db.Close()
		if r.Method == "GET" {
			if err == nil {
				var passengerTrips []Trip = gettrips(db, passengerID)
				if len(passengerTrips) > 0 { //passenger has more than 1 trips
					json.NewEncoder(w).Encode(passengerTrips)
				} else {
					w.WriteHeader(http.StatusNotFound)
				}
			} else {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			}
		}

	}
}

//db function to update trip info when started
func startTripDB(db *sql.DB, DriverID string) {
	query := fmt.Sprintf("UPDATE Trip SET Status = 'In Transit' , PickUpTime = NOW() WHERE DriverID = %s AND Status = 'On The Way'", DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//db function to update trip info when ended
func endTripDB(db *sql.DB, DriverID string) {
	query := fmt.Sprintf("UPDATE Trip SET Status = 'Completed' , DropOffTime = NOW() WHERE DriverID = %s AND Status = 'In Transit'", DriverID)
	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

//function to run a get request to passenger to ensure passenger is not in trip before trip creation
func checkpassenger(passengerID int) bool {
	url := fmt.Sprintf("http://localhost:5001/api/passenger/check/%d", passengerID)
	response, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	if response.StatusCode == http.StatusAccepted {
		println("status accepted in trip")
		return true
	} else {
		println("status failed in trip")
		return false
	}
}

//function to run a post request to driver to update driver's availablity once trip started
func updateDriver(driverID int) {
	url := fmt.Sprintf("http://localhost:5002/api/driver/update/%d", driverID)
	response, err := http.Post(url, "application/json", nil)

	if err != nil {
		panic(err.Error())
	} else {
		response.Body.Close()
	}
}

//function to run a post request to passenger to update passenger's 'intrip' status once trip started
func updatePassenger(passengerID int) {
	url := fmt.Sprintf("http://localhost:5001/api/passenger/%d", passengerID)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		response.Body.Close()
	}
}

//function to get first available driver for trip when created
func getAvailableDriver() Driver {
	fmt.Println("code ran here in trip")
	url := "http://localhost:5002/api/driver"
	response, err := http.Get(url)
	var tripDriver Driver
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if response.StatusCode == http.StatusCreated {
			data, _ := ioutil.ReadAll(response.Body)
			response.Body.Close()
			json.Unmarshal(data, &tripDriver)
		} else {
			fmt.Printf("There are no available drivers now")
			return tripDriver
		}
	}
	return tripDriver
}

//db function for getting last tripid
func getLastID(db *sql.DB) int {
	query1 := "SELECT COUNT(*) FROM Trip"
	query2 := "SELECT TripID FROM Trip ORDER BY TripID DESC LIMIT 1"
	var tripCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {

		results.Scan(&tripCount)
	}
	if tripCount > 0 {
		results, err := db.Query(query2) //Run Query
		var ID int
		if err != nil {
			panic(err.Error())
		}
		if results.Next() {
			results.Scan(&ID)
		}
		return ID
	} else {
		return 0
	}
}

//db function to insert new trip into db
func insertTrip(db *sql.DB, newTrip Trip) {
	ID := (getLastID(db) + 1)
	PID := newTrip.PassengerID
	DID := newTrip.DriverID
	P := newTrip.Pickup
	DO := newTrip.Dropoff
	CLN := newTrip.Carlicenseno
	S := "On The Way"
	query := fmt.Sprintf("INSERT INTO Trip (TripID, DriverID, PassengerID, Pickup, Dropoff, CarLicenseNo, Status) VALUES('%d', '%d', '%d','%d','%d','%s','%s')", ID, DID, PID, P, DO, CLN, S)
	_, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}
}

//http function to accept post request from driver when trip is started
func startTrip(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	params := mux.Vars(r)
	driverID := params["driverid"]
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "PUT" {
			_, err := ioutil.ReadAll(r.Body)
			if err == nil {
				startTripDB(db, driverID)
				w.WriteHeader(http.StatusCreated)
			} else {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			}
		}
	}
}

//http function to accept post request from driver when trip is ended
func endTrip(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	params := mux.Vars(r)
	driverID := params["driverid"]
	passengerID := params["passengerid"]
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "PUT" {
			_, err := ioutil.ReadAll(r.Body)
			if err == nil {
				endTripDB(db, driverID)
				driverIDint, err1 := strconv.Atoi(driverID)
				passengerIDint, err2 := strconv.Atoi(passengerID)
				if err1 == nil {
					updateDriver(driverIDint)
				}
				if err2 == nil {
					updatePassenger(passengerIDint)
				}
				w.WriteHeader(http.StatusCreated)
			} else {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			}
		}
	}
}

//http function to handle post requests to create new trip
func trip(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished executing
	defer db.Close()

	if r.Header.Get("Content-type") == "application/json" {

		// POST is for creating new object
		if r.Method == "POST" {
			var newTrip Trip
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newTrip)
				if checkpassenger(newTrip.PassengerID) { //Passenger does exist
					var driver Driver = getAvailableDriver()
					if driver.DriverID == -1 {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("409 - There are no available drivers at the moment!"))
					} else {
						println("Driver accepted in trip")
						newTrip.DriverID = driver.DriverID
						newTrip.Carlicenseno = driver.Carlicenseno
						insertTrip(db, newTrip)
						updateDriver(newTrip.DriverID)
						updatePassenger(newTrip.PassengerID)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("201 - Trip booked!"))
					}
				} else { //passenger doesnt exist
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte("409 - No Such Passenger!"))
				}

			} else {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter trip detais in JSON format!"))
			}
		}
	}
}

func main() {
	// This is to allow the headers, origins and methods all to access CORS resource sharing
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	router := mux.NewRouter()
	//get driver current trip
	router.HandleFunc("/api/trip/gettrip/{driverid}", getDriverCurrentTrip).Methods("GET")
	//start current trip with driver id supplied
	router.HandleFunc("/api/trip/start/{driverid}", startTrip).Methods("PUT")
	//end current trip with both driver & passenger id supplied
	router.HandleFunc("/api/trip/end/{driverid}/{passengerid}", endTrip).Methods("PUT")
	//get all passenger trips with passengerid supplied
	router.HandleFunc("/api/trip/getpassengertrips/{passengerid}", getpassengertrips).Methods("GET")
	//create new trip
	router.HandleFunc("/api/trip/create", trip).Methods("POST")
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", handlers.CORS(headers, origins, methods)(router)))
}
