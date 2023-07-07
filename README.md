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
	1. localhost:8080/tanks/<tankID>/tank-sensors/waterlevel
 	- The tankID is the deviceid
  	- The endpoint returns the data for the water level sensor in the tank
   	JSON RESPONSE =>
    	[
	  {
	    "id": "5da97e3813474778618e2289",
	    "name": "Water Level Sensor",
	    "modified": "2023-07-07T09:13:59.554Z",
	    "created": "2023-07-07T09:11:02.14Z",
	    "time": "2023-07-07T09:23:42.008Z",
	    "meta": {
	      "kind": "WaterLevel"
	    },
	    "value": 2
	  }
	]

 	2. localhost:8080/tanks/<tankID>/tank-sensors/waterlevel/value
  	- The endpoint returns the most recent value in the water level sensor
   	- The response is helpful in showing the current water level value
   	JSON RESPONSE => Returns the value field which is 2

     	3. localhost:8080/tanks/<tankID>/tank-sensors/waterlevel/values
      	- Returns the historical values stored in the water sensor
       	- The data returned here is helpful in ploting graphs
       	JSON RESPONSE =>
	[	  
	  {
	    "value": 2.2,
	    "time": "2023-07-07T12:23:34+03:00"
	  },
	  {
	    "value": 2.1,
	    "time": "2023-07-07T12:23:37+03:00"
	  },
	  {
	    "value": 2,
	    "time": "2023-07-07T12:23:42+03:00"
	  }
	]
 	
#### Water Temperature endpoints
	- The JSON response for this API endpoint are similar to water level endpoint
 
	1. localhost:8080/tanks/<tankID>/tank-sensors/water-temperature

	2. localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/value

	3. localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/values
 
#### Water quality endpoints
	1. localhost:8080/tanks/<tankID>/tank-sensors/
 	```
	[
	  {
	    "id": "201c85cdbda37",
	    "name": "Water Quality Sensor",
	    "modified": "2023-07-07T09:21:50.772Z",
	    "created": "2023-07-07T09:19:55.884Z",
	    "time": "2023-07-07T09:23:15.76Z",
	    "meta": {
	      "kind": "WaterPollutantSensor"
	    },
	    "value": 911
	  }
	]
 	```
	2. localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/value
 	JSON RESPONSE =>
  	{
	  "tdsValue": 911,
	  "waterQuality": "Poor"
	}
 
	3. localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/values
 	[
	  {
	    "tdsValue": 291,
	    "timestamp": "0001-01-01T00:00:00Z",
	    "waterQuality": "Excellent"
	  },
	  {
	    "tdsValue": 431,
	    "timestamp": "0001-01-01T00:00:00Z",
	    "waterQuality": "Good"
	  },
	  {
	    "tdsValue": 911,
	    "timestamp": "0001-01-01T00:00:00Z",
	    "waterQuality": "Poor"
	  }
	]


### Pump API endpoints
























