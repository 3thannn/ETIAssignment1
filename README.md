<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#usage">Usage</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#license">License</a></li>
    
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

**Hello There! Welcome to GOGO Gojek's README page.**

GOGO Gojek is a simple ride hailing platform made with simple HTML, JS, JQUERY, SQL and GOLANG

## Screenshots
**Main Landing Page**

![LandingPage_Screenshot][home-page]
---

**Create Passenger Account**

![Create_Passenger_ScreenShot][create-passenger]
---

**Create Driver Account**

![Create_Driver_ScreenShot][create-driver]
---

**Passenger Login Page**

![Passenger_Login_Screenshot][passenger-login]
---

**Passenger Logged In Page**

![Passenger_Loggedin_Screenshot][passenger-loggedin]
---

**Create New Trip**

![Create_New_Trip][create-trip]
---

**Update Passenger Page**

![Update_Passenger_Screenshot][update-passenger]
---

**Driver Login Page**

![Driver_Login_Screenshot][driver-login]
---

**Driver Logged In Page (In Trip)**

![Driver_Loggedin_No_Trip_Screenshot][Driver-loggedinnotrip]
---

**Driver Logged In Page (In Trip)**

![Driver_Loggedin_In_Trip_Screenshot][Driver-loggedintrip]
---

**Update Driver Page**

![Update_Driver_Screenshot][update-driver]
---

### Built With

* [GOLANG](https://go.dev/)
* [HTML](https://html.com/)
* [JavaScript](https://reactjs.org/)
* [Bootstrap](https://getbootstrap.com)
* [JQuery](https://jquery.com)
* [MYSQL](https://www.mysql.com/)


### Assignment Objectives: 
•	To demonstrate ability to develop REST APIs
•	To make conscientious consideration in designing microservices

### Assignment Requirements: 
- [x] Minimum 2 microservices using Go
- [x]	Persistent storage of information using database, e.g. MySQL

### Bonus Marks: 
- [x] Implement a web frontend that calls your microservices, instead of a console application. You can implement the frontend with any language and design of your choice.

### Microservice Architecture: 


**Project Architecture**

![microservice-diagran][microservice-diagram]
---

There are a total of 3 Microservices used in this application, along with the use of monolith frontend.

1. **Passenger**<br />
  Functions consists of:
    - Create Passenger
    - Update Driver
    - Create Trip 
      - Calls Trip Microservice
2. **Driver**<br />
  Functions consists of:
    - Create Driver
    - Update Driver
    - Start Trip
      - Calls Trip Microservice
   - End Trip
      - Calls Trip Microservice
3. **Trip**<br />
  Functions consists of:
    - Create Trip (Passenger Creates Trip) 
      - Calls Passenger & Driver Mircoservice to run validation
    - Get All Passenger Trips
    - Start Trip (Driver Starts Trip)
    - End Trip (Driver Ends Trip)
<!-- GETTING STARTED -->
### Getting Started
Below are the steps needed to get GOGO GOJEK running on your own system!

### Prerequisites<br />

  Please ensure that GOLANG and MYSQL is installed on your system, and is fully operational

  Please do also ensure that your SQL user login is as such:
 ```sh
    Username: user
    Password: password
 ```
### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/3thannn/ETIAssignment1.git
   ```
2. Install necessary libraries
   ```sh
   go get -u github.com/go-sql-driver/mysql
   go get -u github.com/gorilla/mux
   go get -u github.com/gorilla/handlers
   ```
3. Execute database script located in `/db/InitDB.sql`



<!-- USAGE EXAMPLES -->
## Usage


To start the platform, please follow the following steps:
1. Run all the respective microservices in your system
   ```sh
    cd ETIAssignment1/Trip
    go run Trip.go
   ```

    ```sh
      cd ETIAssignment1/Driver
      go run Driver.go
     ```
    ```sh
      cd ETIAssignment1/Trip
      go run Trip.go
     ```

2. Open webpage by opening `index.html` in `ETIAssignment1/FrontEnd/Main`

<!-- ROADMAP -->
## Roadmap

- [x] Create Database
- [x] Create Microservices Backend using GOLANG
- [x] Create Frontend pages using HTML








<!-- CONTACT -->
## Contact

My Email - [School Email](mailto:s10185214@connect.np.edu.sg) 

<p align="left">(<a href="#top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->

[home-page]: ./Images/Home_Page.PNG
[create-driver]: ./Images/Create_Driver.PNG
[create-passenger]: ./Images/Create_Passenger.PNG
[passenger-login]: ./Images/Passenger_Login.PNG
[passenger-loggedin]: ./Images/Passenger_LoggedIn.PNG
[create-trip]: ./Images/Create_NewTrip
[update-passenger]: ./Images/Update_Passenger.PNG
[driver-login]: ./Images/Driver_Login.PNG
[driver-loggedintrip]: ./Images/DriverLogin_InTrip.PNG
[driver-loggedinnotrip]: ./Images/DriverLogin_NoTrip.PNG
[update-driver]: ./Images/Update_Driver.PNG
[microservice-diagram]: ./Images/Microservice_Diagram.PNG


