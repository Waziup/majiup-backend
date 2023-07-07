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
### Tank API endpoints
	1. GET = localhost:8080/tanks
 		Returns the tanks registered in the gateway	
 
	2. GET = localhost/tanks/<tankID>
 		Returns the specific tank with the given tank id
   
    	3. GET = localhost/tanks/<tankID/tank-sensors
     		Returns the sensors that are connected to the tank       	

### Sensors API endpoints
#### Water Level endpoints
	// Endpoint to get the water level sensor data from a specific tank
	1. GET = localhost:8080/tanks/<tankID>/tank-sensors/waterlevel
 	- The tankID is the deviceid
  	- The endpoint returns the data for the water level sensor in the tank
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

 	2. GET = localhost:8080/tanks/<tankID>/tank-sensors/waterlevel/value
  	- The endpoint returns the most recent value in the water level sensor
   	- The response is helpful in showing the current water level value
   	JSON RESPONSE => Returns the value field which is **2**

     	3. GET = localhost:8080/tanks/<tankID>/tank-sensors/waterlevel/values
      	- Returns the historical values stored in the water sensor
       	- The data returned here is helpful in ploting graphs
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
 
	1. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature

	2. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/value

	3. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/values
 
#### Water quality endpoints
	1. GET = localhost:8080/tanks/<tankID>/tank-sensors/
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
 
	2. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/value
  	{
	  "tdsValue": **911**,
	  "waterQuality": "**Poor**"
	}
 
	3. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/values
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

	1. GET = localhost:8080/tanks/<tankID>/pumps
 		- Returns the pump that is specific to the given tank
   		[
		  {
		    "id": "201c85cdbda3",
		    "name": "Pump",
		    "modified": "2023-07-07T09:28:44.035Z",
		    "created": "2023-07-07T09:27:19.222Z",
		    "meta": {
		      "kind": "Motor"
		    },
		    "time": "2023-07-07T09:58:07.36Z",
		    "value": 1
		  }
		]


	2. GET = localhost:8080/tanks/<tankID>/pumps/state
 		- Returns the pump state as a **1** or **0**

	3. GET = localhost:8080/tanks/<tankID>/pumps/states
 		- Returns the historical state of the pump
   		- The length of the json array is equivalent to how many times the pump has been actuated
     		-[
		  {
		    "pumpState": 1,
		    "timestamp": "0001-01-01T00:00:00Z"
		  },
		  {
		    "pumpState": 0,
		    "timestamp": "0001-01-01T00:00:00Z"
		  },
		  {
		    "pumpState": 1,
		    "timestamp": "0001-01-01T00:00:00Z"
		  },		 
		]

	4. POST = localhost:8080/tanks/<tankID>/pumps/state
 		curl
   		curl -X POST -H "Content-Type: application/json" -d '{"value": 1}' 	http://localhost:8080/tanks/201c85cdbda3/pumps/state
     
 
























