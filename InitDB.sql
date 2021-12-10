	DROP DATABASE ETIAssignment1;

	CREATE DATABASE ETIAssignment1;

	USE ETIAssignment1;

	DROP TABLE IF EXISTS Passenger;
	CREATE TABLE Passenger(
		PassengerID int NOT NULL PRIMARY KEY,
		FirstName varchar(255) NOT NULL,
		LastName varchar(255) NOT NULL,
		MobileNo varchar(8) NOT NULL,
		EmailAddress varchar(255) NOT NULL,
        InTrip bool DEFAULT False ,
		UNIQUE KEY unique_email (EmailAddress),
		UNIQUE KEY unique_mobile (MobileNo));


	DROP TABLE IF EXISTS Driver;
	CREATE TABLE Driver(
		DriverID int NOT NULL PRIMARY KEY,
		FirstName varchar(255) NOT NULL,
		LastName varchar(255) NOT NULL,
		MobileNo varchar(8) NOT NULL,
		EmailAddress varchar(255) NOT NULL,
        NRIC varchar(9) NOT NULL,
		CarLicenseNo varchar(8) NOT NULL,
		Available bool NOT NULL DEFAULT True,
		UNIQUE KEY unique_email (EmailAddress) ,
		UNIQUE KEY unique_mobile (MobileNo),
		UNIQUE KEY unique_liscno (CarLicenseNo));



	DROP TABLE IF EXISTS Trip;
	CREATE TABLE Trip(
		TripID int NOT NULL PRIMARY KEY,
		DriverID int NOT NULL,
		PassengerID int NOT NULL,
		Pickup int(6) NOT NULL,
		Dropoff int(6) NOT NULL,
		PickUpTime datetime,
		DropOffTime datetime,
        CarLicenseNo varchar(8) NOT NULL,
		Status varchar(255) NOT NULL,
		CONSTRAINT CHK_Status CHECK (Status IN ('On The Way', 'In Transit', 'Completed')));
	










