# majiup-backend

## Steps to run the Majiup Backend Locally on your computer

### Requirements
Golang - Install the latest version of GO

### Step 1 - Clone the repository
```
git clone https://github.com/JosephMusya/majiup-backend.git
```
### Step 2 - Navigate to the majiup-backend repository
```
cd majiup-backend
```
### Step 2 - Download the go dependencies
```
go mod download
```
### Step 4 - Build the binary file
```
go build -o main main.go
```
### Step 5 - Run the binary file
```
sudo ./main
```
You can alternatively run the main.go file
```
go run main.go
```
## Steps to run the Majiup Backend in docker
//

## Available API endpoints
The majiup-backend acts as a proxy api to http://localhost/devices which is wazigate api endpoint

### Sensors API endpoints
#### Water Level endpoints
	// Endpoint to get the water level sensor data from a specific tank
	localhost:8080/tanks/<tankID>/tank-sensors/waterlevel
 

	// Endpoint to get the water level value
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/value", GetWaterLevelValueHandler)

	// Endpoint to get the water level history values
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/values", GetWaterLevelHistoryHandler)

### Pump API endpoints
























