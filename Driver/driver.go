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

//db function to get driver object using email address
func getDriverLogin(db *sql.DB, emailAddr string) Driver {
	query := fmt.Sprintf("SELECT * FROM Driver WHERE EmailAddress = '%s'", emailAddr)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var driver Driver
	for results.Next() {
		results.Scan(&driver.DriverID, &driver.Firstname, &driver.Lastname, &driver.Mobileno, &driver.Emailaddress, &driver.Nric, &driver.Carlicenseno, &driver.Availability)
	}
	return driver
}

//db function to update driver's availablility
func updateAvailablility(db *sql.DB, driverID string) {
	ID := driverID
	query1 := fmt.Sprintf("SELECT Available FROM Driver WHERE DriverID = %s", ID)
	query2 := fmt.Sprintf("UPDATE Driver SET Available = false WHERE DriverID = %s", ID)
	query3 := fmt.Sprintf("UPDATE Driver SET Available = true WHERE DriverID = %s", ID)

	results, err := db.Query(query1)
	if err != nil {
		panic(err.Error())
	}
	var available bool
	if results.Next() {
		results.Scan(&available)
	}
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

//db function to check if there is any available drivers
func anyavailabledriver(db *sql.DB) bool {
	println("code ran here in anyavailabledriver")
	query := "SELECT COUNT(*) FROM Driver WHERE Available = 1"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		var count int
		fmt.Println(&count)
		results.Scan(&count)
		if count > 0 {
			return true //return true if at least 1 available driver exists
		}
	}
	return false
}

//db function to check if a driver object already exists with entered email and eddress
func checkdriverexists(db *sql.DB, newDriver Driver) bool {
	MN := newDriver.Mobileno
	EA := newDriver.Emailaddress
	NRIC := newDriver.Nric
	CLN := newDriver.Carlicenseno
	query := fmt.Sprintf("SELECT COUNT(*) FROM Driver WHERE MobileNo = '%s' OR EmailAddress = '%s' OR NRIC ='%s' OR Carlicenseno = '%s'", MN, EA, NRIC, CLN)

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

//db function to get first available driver to assign to trip
func getfirstavailabledriver(db *sql.DB) Driver {
	query := "SELECT * FROM Driver WHERE Available = 1 ORDER BY DriverID ASC LIMIT 1;"
	var lastDriver Driver
	//check if there is indeed any available drivers
	results, err := db.Query(query) //Run Query
	if err != nil {
		panic(err.Error())
	}

	if results.Next() {
		results.Scan(&lastDriver.DriverID, &lastDriver.Firstname, &lastDriver.Lastname, &lastDriver.Mobileno, &lastDriver.Emailaddress, &lastDriver.Nric, &lastDriver.Carlicenseno, &lastDriver.Availability)
	}
	fmt.Println(lastDriver)
	return lastDriver //return driver object

}

//db function to get last driverid
func getlastid(db *sql.DB) int { //Gets the last id of driver
	query1 := "SELECT COUNT(*) FROM Driver"
	query2 := "SELECT DriverID FROM Driver ORDER BY DriverID DESC LIMIT 1"
	var driverCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {

		results.Scan(&driverCount)
	}
	if driverCount > 0 {
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

//db function to insert new created driver account
func insertRecord(db *sql.DB, newDriver Driver) {
	ID := (getlastid(db) + 1)
	FN := newDriver.Firstname
	LN := newDriver.Lastname
	MN := newDriver.Mobileno
	EA := newDriver.Emailaddress
	NRIC := newDriver.Nric
	CLN := newDriver.Carlicenseno
	A := 1
	query := fmt.Sprintf("INSERT INTO Driver VALUES ('%d', '%s', '%s', '%s', '%s', '%s', '%s', '%d')",
		ID, FN, LN, MN, EA, NRIC, CLN, A) //Create Query to insert driver into db

	_, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}
}

//db function to update driver details
func updateRecord(db *sql.DB, newDriver Driver) {
	ID := newDriver.DriverID
	FN := newDriver.Firstname
	LN := newDriver.Lastname
	MN := newDriver.Mobileno
	EA := newDriver.Emailaddress
	CLN := newDriver.Carlicenseno
	query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNo = '%s', EmailAddress = '%s' , Carlicenseno = '%s' WHERE DriverID = '%d' ",
		FN, LN, MN, EA, CLN, ID) //Create Query to update Driver in db
	_, err := db.Query(query) //Run Query

	if err != nil {
		panic(err.Error())
	}
}

//http function to handle get request for driver details using email login
func loginDriver(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")
	if err != nil {
		panic(err.Error())
	}
	params := mux.Vars(r)
	value := params["emailaddress"]

	if r.Method == "GET" {
		if err == nil {
			if getDriverLogin(db, value).Emailaddress != "" {
				json.NewEncoder(w).Encode(getDriverLogin(db, value))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		}
	}
}

//http function to handle GET, POST and PUT requests.
//GET -> Get driver details for trip
//POST -> Driver account creation
//PUT -> Update driver details

func driver(w http.ResponseWriter, r *http.Request) {
	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")

	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	if r.Method == "GET" {
		println("code ran here in driver get")
		if anyavailabledriver(db) {
			tripDriver := getfirstavailabledriver(db)
			w.WriteHeader(http.StatusCreated)
			fmt.Println("201 - Driver found!")
			json.NewEncoder(w).Encode(tripDriver)
		} else {
			fmt.Println("409 - No Available Driver found!")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Nil"))
			tripDriver := Driver{DriverID: -1}
			json.NewEncoder(w).Encode(tripDriver)
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new object
		if r.Method == "POST" {

			// read the string sent to the service
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)
				// check if course exists; add only if
				// course does not exist
				if !checkdriverexists(db, newDriver) {
					insertRecord(db, newDriver)
					w.WriteHeader(
						http.StatusCreated)
					w.Write([]byte("201 - Account successfully created!"))
				} else {
					w.WriteHeader(
						http.StatusConflict)
					w.Write([]byte("409 - Account with details already exists!"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please enter account detais in JSON format!"))
			}
		}
		if r.Method == "PUT" {
			// read the string sent to the service
			var driverInput Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &driverInput)
				updateRecord(db, driverInput)
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("202 - Driver ID:" + strconv.Itoa(driverInput.DriverID) + "has been successfully updated!"))
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Invalid JSON Format!"))
			}
		}
	}
}

//http function for post request, update driver's availablity
func updateDriver(w http.ResponseWriter, r *http.Request) {
	// Use mysql as driverName and a valid DSN as dataSourceName:
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ETIAssignment1")

	// handle error
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			params := mux.Vars(r)
			driverID := params["driverid"]

			if r.Method == "POST" {
				if err == nil {
					updateAvailablility(db, driverID)
				}
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
	//handle driver login
	router.HandleFunc("/api/driver/login/{emailaddress}", loginDriver).Methods("GET")
	//handle driver get for trip
	//handle post for driver account creation
	//handle put for driver details update
	router.HandleFunc("/api/driver", driver).Methods("GET", "POST", "PUT")
	//handle post request for updating driver availability
	router.HandleFunc("/api/driver/update/{driverid}", updateDriver).Methods("POST")
	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", handlers.CORS(headers, origins, methods)(router)))
}
