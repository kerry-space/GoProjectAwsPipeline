# Go API Carsystem management

system, which offers a RESTful API for creating, reading, updating, and deleting cars (CRUD) as well as managing user registrations and sessions. The application is built with Go and Gin framework and uses MySQL for database management and Redis for session management.

## Installation

To install and use this API, follow these steps:

1. Clone the repository:
   ```sh
   git clone 
2. Install dependencies:
    ```sh
   go mod tidy
3. Build and run the API
   Ensure you have MySQL and Redis installed and running on your system
   ```sh
   go run main.go

# Docker Containerization

1. Build Docker container
   ```sh
   docker build -t car-api.
2. Run Docker container:
   ```sh
   docker run -p 8080:8080 -d car-api

# Docker MySQL and Redis Setup

1.  Pull MySQL Docker image:
   ```sh
   docker pull mysql
2. Run MySQL Docker container:
   ```sh
   docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 -d mysql
3. Pull Redis Docker image:
  ```sh
  docker pull redis
3. Run Redis Docker container:
   ```sh
   docker run --name redis-container -p 6379:6379 -d redis
