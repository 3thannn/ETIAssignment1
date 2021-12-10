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

type Passenger struct {
	PassengerID int

	Firstname    string
	Lastname     string
	Mobileno     string
	Emailaddress string
	Intrip       bool
}

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

// func validKey(r *http.Request) bool {
// 	v := r.URL.Query()
// 	if key, ok := v["key"]; ok {
// 		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
// 			return true
// 		} else {
// 			return false
// 		}
// 	} else {
// 		return false
// 	}
// }

// func home(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to the REST API!")
// }
func getPassenger(db *sql.DB, emailAddr string) Passenger {
	println(emailAddr)
	query := fmt.Sprintf("SELECT * FROM Passenger WHERE EmailAddress = '%s'", emailAddr)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var passenger Passenger
	for results.Next() {
		results.Scan(&passenger.PassengerID, &passenger.Firstname, &passenger.Lastname, &passenger.Mobileno, &passenger.Emailaddress, &passenger.Intrip)
	}
	return passenger
}

//function to call a get request from trip to get all customer trips back to be used in passenger context
func getTrips(passengerID string) []Trip {
	url := fmt.Sprintf("http://localhost:5003/api/trip/getpassengertrips/%s", passengerID)
	response, err := http.Get(url)
	var passengerTrips []Trip
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if response.StatusCode == http.StatusNotFound {
			fmt.Println("409 - No trips found!")
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			response.Body.Close()
			json.Unmarshal(data, &passengerTrips)
			fmt.Println("202 - Successfully received trips!")
		}
	}
	return passengerTrips
}

//db function to update passenger's status to intrip
func updateInTrip(db *sql.DB, passengerID string) {
	ID := passengerID
	query1 := fmt.Sprintf("SELECT InTrip FROM Passenger WHERE PassengerID = %s", ID)
	query2 := fmt.Sprintf("UPDATE Passenger SET InTrip = false WHERE PassengerID = %s", ID)
	query3 := fmt.Sprintf("UPDATE Passenger SET InTrip = true WHERE PassengerID = %s", ID)

	results, err := db.Query(query1)
	if err != nil {
		panic(err.Error())
	}
	var available bool
	if results.Next() {
		results.Scan(&available)
	}
	fmt.Println(available)
	if available {
		_, err := db.Query(query2) //Run Query

		if err != nil {
			panic(err.Error())
		}
	} else {
		_, err := db.Query(query3) //Run Query

		if err != nil {
			panic(err.Error())
		}
	}
}

//db function to insert new passenger object into db after validated
func insertRecord(db *sql.DB, newPassenger Passenger) {
	ID := (getlastid(db) + 1)
	FN := newPassenger.Firstname
	LN := newPassenger.Lastname
	MN := newPassenger.Mobileno
	EA := newPassenger.Emailaddress
	query := fmt.Sprintf("INSERT INTO Passenger (PassengerID, FirstName, LastName, MobileNo, EmailAddress) VALUES ('%d', '%s', '%s', '%s', '%s')",
		ID, FN, LN, MN, EA) //Create Query to insert passenger into db

	_, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}
}

//db function to update passenger object details inside
func updateRecord(db *sql.DB, newPassenger Passenger) {
	ID := newPassenger.PassengerID
	FN := newPassenger.Firstname
	LN := newPassenger.Lastname
	MN := newPassenger.Mobileno
	EA := newPassenger.Emailaddress
	query := fmt.Sprintf("UPDATE Passenger SET FirstName = '%s', LastName = '%s', MobileNo = '%s', EmailAddress = '%s' WHERE PassengerID = '%d' ",
		FN, LN, MN, EA, ID) //Create Query to insert passenger into db

	_, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}
}

//gets the last passengerid from db
func getlastid(db *sql.DB) int { //Gets the last id of passengers
	query1 := "SELECT COUNT(*) FROM Passenger"
	query2 := "SELECT PassengerID FROM Passenger ORDER BY PassengerID DESC LIMIT 1"
	var passengerCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		results.Scan(&passengerCount)
	}
	if passengerCount > 0 {
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

//This function checks that passenger in question isn't in a trip
//Returns true if in trip, false if not in trip
func checkpassenger(db *sql.DB, passengerID string) bool {
	query := fmt.Sprintf("SELECT COUNT(*) FROM Passenger WHERE PassengerID = '%s' AND InTrip = 0", passengerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var count int
	if results.Next() {
		results.Scan(&count)
	}
	if count != 0 {
		return true
	} else {
		return false
	}
}

//html function for trip to check if passenger is in trip or not
//returns status accepted if passed
//returns status not found if failed
func checkpassengerexists(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	}
	params := mux.Vars(r)
	passengerID := params["passengerid"]
	if r.Method == "GET" {
		defer db.Close()
		if checkpassenger(db, passengerID) {
			fmt.Println("Checkpassenger not intrip passed")
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Passenger ID:" + passengerID + "exists!"))
		} else {
			fmt.Println("Checkpassenger not intrip failed")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("409 - No such passenger!"))
		}
	}
}

//checks if passenger with mobile number or email addres already exists
//validation for creation of passenger account
func doespassengerexists(db *sql.DB, newPassenger Passenger) bool {
	MN := newPassenger.Mobileno
	EA := newPassenger.Emailaddress
	query := fmt.Sprintf("SELECT COUNT(*) FROM PASSENGER WHERE MobileNo = '%s' or EmailAddress = '%s'", MN, EA) //Create Query to insert passenger into db

	results, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		var Count int
		results.Scan(&Count)
		if Count > 0 {
			return true
		}
	}
	return false
}

//http function for getting all trips
func getAllTrips(w http.ResponseWriter, r *http.Request) {
	//GET is for getting passenger trips
	if r.Method == "GET" {
		params := mux.Vars(r)
		passengerID := params["passengerid"]
		var passengerTrips []Trip = getTrips(passengerID)
		json.NewEncoder(w).Encode(passengerTrips)
		fmt.Println(passengerTrips)
	}
}

//http function to handle the creation of customer through POST request
func createPassenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")

	// handle error
	if err != nil {
		panic(err.Error())
	}
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if !doespassengerexists(db, newPassenger) {
					insertRecord(db, newPassenger)
					w.WriteHeader(
						http.StatusCreated)
					w.Write([]byte("201 - Account successfully created!"))
				} else {
					w.WriteHeader(
						http.StatusConflict)
					w.Write([]byte("409 - Account with entered details already exists!"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account detais in JSON format!"))
			}
		}
	}
}

//http function to handle
//GET request for passenger object
//POST request for updating passenger object status
//PUT request for updating passenger object details
func passenger(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")

	// handle error
	if err != nil {
		panic(err.Error())
	}
	params := mux.Vars(r)
	value := params["value"]

	if r.Method == "GET" {
		if err == nil {
			if getPassenger(db, value).Emailaddress != "" {
				fmt.Println(getPassenger(db, value))
				json.NewEncoder(w).Encode(getPassenger(db, value))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
	}

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new passenger

		if r.Method == "POST" {
			if err == nil {
				updateInTrip(db, value)
			}
		}
		if r.Method == "PUT" {
			// read the string sent to the service
			var passengerInput Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &passengerInput)
				if !doespassengerexists(db, passengerInput) {
					insertRecord(db, passengerInput)
					w.WriteHeader(
						http.StatusCreated)
					w.Write([]byte("201 - Account successfully updated!"))
				} else {
					updateRecord(db, passengerInput)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Passenger ID:" + strconv.Itoa(passengerInput.PassengerID) + "has been successfully updated!"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account detais in JSON format!"))
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
	//check customer exists
	router.HandleFunc("/api/passenger/check/{passengerid}", checkpassengerexists).Methods("GET")
	//create new passenger
	router.HandleFunc("/api/passenger/create", createPassenger).Methods("POST")
	//get passenger trips
	router.HandleFunc("/api/passenger/trips/{passengerid}", getAllTrips).Methods("GET")
	// get passenger object (GET)
	// update passenger details (PUT)
	// update passenger status (POST)
	router.HandleFunc("/api/passenger/{value}", passenger).Methods("GET", "PUT", "POST")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", handlers.CORS(headers, origins, methods)(router)))
}
