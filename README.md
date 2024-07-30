# Golang Excel Data Importer

This project is a Golang-based application that imports data from an Excel file, stores it into a MySQL database, and caches the data in Redis. It also provides a simple CRUD (Create, Read, Update, Delete) system to manage the imported data.

## Features

- Import data from an Excel file (.xlsx)
- Store data in MySQL
- Cache data in Redis
- CRUD operations for managing data

## Technologies Used

- [Golang](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Excelize](https://github.com/xuri/excelize)
- [MySQL](https://www.mysql.com/)
- [Redis](https://redis.io/)

## Prerequisites

- Go 1.16+
- MySQL
- Redis

## Getting Started

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/golang-excel-importer.git
    cd golang-excel-importer
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

### Configuration

1. **Create a `.env` file:**

    Create a `.env` file in the root of your project with the following content:

    ```env
    DB_URI=username:password@tcp(localhost:3306)/dbname?parseTime=True&loc=Local
    REDIS_URI=redis://:@localhost:6379/0
    ```

    - `DB_URI`: This variable contains the connection string for your MySQL database. Replace `root`, `123456789`, `localhost`, `3306`, and `importdb` with your MySQL username, password, host, port, and database name respectively.
    - `REDIS_URI`: This variable contains the connection URI for your Redis instance. Adjust the URI according to your Redis configuration.

2. **Database Configuration:**

    The application will automatically use the `DB_URI` from the `.env` file to connect to the MySQL database.

3. **Redis Configuration:**

    The application will automatically use the `REDIS_URI` from the `.env` file to connect to the Redis instance.

### Run the Application

1. **Start the server:**

    ```bash
    go run main.go
    ```

2. **Access the API:**

    The server will be running on `http://localhost:8080`.

## API Endpoints

### Upload Excel File

- **URL:** `/upload`
- **Method:** `POST`
- **Description:** Upload an Excel file to import data.
- **Parameters:**
  - `file`: Form file parameter for the `.xlsx` file.
  - It only supports for `.xlsx` file .
- **Curl:** 
curl --location 'http://127.0.0.1:8080/upload' \
--form 'file=@"postman-cloud:///1ef4db57-2d5d-4390-babf-14353d2f52ee"'


### Get Records

- **URL:** `/records?limit=10&offset=10`
- **Method:** `GET`
- **Description:** Retrieve all records from the database with pagination.
- **Parameters:**
  - `limt`: Query parameter for the limit.
  - `offset`: Query parameter for the offset.
- **Curl:** 
curl --location 'http://127.0.0.1:8080/records?limit=3&offset=3'

### Get Single Record

- **URL:** `/records/:id`
- **Method:** `GET`
- **Description:** Get a record by ID.
- **Parameters:**
  - `id`: Path parameter for the record ID.
- **Curl:** 
curl --location 'http://127.0.0.1:8080/records/1582'

### Update Record

- **URL:** `/records/:id`
- **Method:** `PUT`
- **Description:** Update a record by ID.
- **Parameters:**
  - `id`: Path parameter for the record ID.
  - JSON body with updated record details.
- **Curl:** 
curl --location --request PUT 'http://127.0.0.1:8080/records/1582' \
--header 'Content-Type: application/json' \
--data '{
    "record_id":1582,
    "first_name":"Dnyanesh",
    "last_name":"Rode"
}'

### Delete Record

- **URL:** `/records/:id`
- **Method:** `DELETE`
- **Description:** Delete a record by ID.
- **Parameters:**
  - `id`: Path parameter for the record ID.
- **Curl:** 
curl --location --request DELETE 'http://127.0.0.1:8080/records/1562'

## Example JSON for Update Record

```json
{
    "record_id":1582,
    "first_name":"Dnyanesh",
    "last_name":"Rode"
}
