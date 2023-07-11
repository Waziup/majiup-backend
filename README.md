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

	4. GET = localhost/tanks/<tankID>/tank-info
		Returns the 3 sensor history data with there timestamps

		{
			"waterLevel": [				
				{
				"timestamp": "2023-07-07T09:23:42.008Z",
				"value": 2.2
				},
				{
				"timestamp": "2023-07-07T09:23:42.008Z",
				"value": 2.1
				},
				{
				"timestamp": "2023-07-07T09:23:42.008Z",
				"value": 2
				}
			],
			"waterTemperature": [				
				{
				"timestamp": "2023-07-07T09:24:42.625Z",
				"value": 24.5
				},
				{
				"timestamp": "2023-07-07T09:24:42.625Z",
				"value": 23
				},
				{
				"timestamp": "2023-07-07T09:24:42.625Z",
				"value": 23.3
				}
			],
			"waterQuality": [				
				{
				"timestamp": "2023-07-07T09:23:15.76Z",
				"value": 291
				},
				{
				"timestamp": "2023-07-07T09:23:15.76Z",
				"value": 431
				},
				{
				"timestamp": "2023-07-07T09:23:15.76Z",
				"value": 911
				}
			]
		}


### Sensors API endpoints
#### Water Level endpoints
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
 	- Sensor kind is "WaterThermometer"	
 
	1. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature
 		[
		  {
		    "id": "5da97e3813474778618e2289",
		    "name": "Water Temperature Sensor",
		    "modified": "2023-07-07T09:13:59.554Z",
		    "created": "2023-07-07T09:11:02.14Z",
		    "time": "2023-07-07T09:23:42.008Z",
		    "meta": {
		      "kind": "WaterThermometer"
		    },
		    "value": 23.6
		  }
		]

	2. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/value
 		- Returns  a singe value, 23.6

	3. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-temperature/values
 		- Returns a list of sensor values collected by the sensor and their respective timestamps
 
#### Water quality endpoints
	1. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-quality
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
 
	2. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-quality/value
	  	{
		  "tdsValue": **911**,
		  "waterQuality": "**Poor**"
		}
	 
	3. GET = localhost:8080/tanks/<tankID>/tank-sensors/water-quality/values
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
     
 ### Tank meta ( Settings, Notifications, Location) API endpoints
	-The meta is a field within the tank and can be update using the endpoint,
 	GET, POST = localhost/devices/<device-id>/meta
	Meta Structure:
		"meta": {
		    "receivenotifications": false,
		    "notifications": {
		      "id": "",
		      "message": "",
		      "read_status": false
		    },
		    "location": {
		      "longitude": 0,
		      "latitude": 0
		    },
		    "settings": {
		      "length": 0,
		      "width": 0,
		      "height": 0,
		      "radius": 0,
		      "capacity": 0,
		      "maxalert": 0,
		      "minalert": 0
		    }
		  },

 	1. POST = localhost:8080/tanks/:tankID/location => To post location cordincates
  		- curl
    		curl -X POST -H "Content-Type: application/json" -d '{"latitude": 9, "longitude": 32.438570}' http://localhost:8080/tanks/201c85cdbda3/location

    ********** **To be completed (APIs)** *************
    	localhost:8080/tanks/:tankID/location => To return the location cordinates
	
	GET = localhost:8080/tanks/:tankID/settings => To get the tank settings
	POST = localhost:8080/tanks/:tankID/settings => To post tank settings
	
	GET = localhost:8080/tanks/notifications => To get all tank notifications
	localhost:8080/tanks/:tankID/notifications => To get specific tank notifications
	POST = localhost:8080/tanks/:tankID/notifications => To post a tank notification
	DELETE = localhost:8080/tanks/:tankID/notification => To delete a tank notification
 
























