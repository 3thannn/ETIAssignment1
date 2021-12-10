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

Hello There!
Welcome to GOGO Gojek's README page.

GOGO Gojek is a simple ride hailing platform made with simple HTML, JS, JQUERY, SQL and GOLANG



### Built With

This section should list any major frameworks/libraries used to bootstrap your project. Leave any add-ons/plugins for the acknowledgements section. Here are a few examples.

* [GOLANG](https://go.dev/)
* [HTML](https://html.com/)
* [JavaScript](https://reactjs.org/)
* [Bootstrap](https://getbootstrap.com)
* [JQuery](https://jquery.com)
* [MYSQL](https://www.mysql.com/)





<!-- GETTING STARTED -->
## Getting Started
Below are the steps needed to get GOGO GOJEK running on your own system!

### Prerequisites

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

