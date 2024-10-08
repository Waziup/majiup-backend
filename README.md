# Majiup backend

## Steps to run the Majiup

The frontend is served from the backend, therefore, there is no need to download the frontend files. <br />
Note that the application is running on the wazigate ip address, to access the application you must have your gateway powered.<br />
The application can be accessed from https://wazigate.local:8081<br />
**In place of wazigate.local, you can key in the ip address of your gateway followed by the port number the app is running on.**

Example http://192.168.0.104:8081

### Step 0 - SSH to gateway

```
ssh pi@wazigate.local
```

Clone the Majiup repository

```
git clone https://github.com/Waziup/majiup-backend.git
```

Navigate into the repository

```
cd majiup-backend
```

### Step 1 - Pull the Majiup Image from dockerhub

```
docker pull waziupiot/majiup:v1.1
```

Confirm that the image is pulled successfully

```
docker images
```

The Majiup image should be among the images

### Step 2 - Create the docker container from the docker-compose file

```
docker-compose up -d
```

The container is build and run in detached mode

### Step 3 - Check the container ID

```
docker ps -a
```

Confirm the container ID

### Step 4 - Container trace logs

You can check the container logs by running the following command

```
sudo docker logs <majiup-container>
```

The api is served by http://wazigate.local:8081/api/v1/

## Creating your Majiup device

### Step 1 - Open the wazigate UI

Open the gateway UI dashboard from http://wazigate.local and navigate to the dashbaord

Create a new device and assign unique adresses to the device.

Make the device LoRAWAN and allocate XLPP for data transmission.

On setting up the device on the gateway. Proceed to setup the hardware and wait for the sensors to be allocated automatically when the hardware on the tank sends data.

### Step 2 - Setup your hardware

Upload this https://github.com/Waziup/majiup-hardware/tree/main/majiup-hardware code to your hardware. Normally, the sensor sends data at inteval of 5 minutes.

Note: Edit the sensor pins defined in the code and the device address as connected to your hardware before uploading.

<!-- e.g In the code, the TDS sensor pin is A1, temperature probe connected to pin A2, echo pin and trigger pins are D3 and D4 respectively -->

When done changing the necessary fields, upload your code and set the device ready to make measurements. Turn the device on when you are done mounting to avoid streaming false value. Even though, this has been catered in the code to reject sending outlier values (false data).

### Step 3 - Set sensors on the gateway

After receiving sensor values on the gateway, edit their Kinds respectively

Normall,the sensors will have name like temperature sensor 1,....

Sensor with ID of temperature_sensor_0 is water level sensor -> Assign kind to WaterLevel

<!-- Sensor with ID of temperature_sensor_1 is water temperature sensor -> Assign kind to WaterThermometer

Sensor with ID of temperature_sensor_2 is water quality sensor sensor -> Assign kind to WaterPollutant -->

### Step 4 - Set your tank on Majiup application

Head over to majiup application and set the dimensions and capacity of your tank under settings.

# API DOCUMENTATION

- The base url for the API is `localhost:<PORT>/api/v1/`

1. Retrieving and posting gateway information
   This endpoint takes bothe POST and GET methods
   - `/gateway-profile`
2. Sending notification
   - `/send-notification`
3. Retrieving battery info
   - `/tanks/{tankID}/battery-info`
4. Retrieving analytics from a particular tank
   - `/tanks/{tankID}/analytics`
5. Listing all tanks connected to the gateway
   - `/tanks`
6. Retrieving a particular tank by ID
   This endpoint takes Get, Delete,
   - `/tanks/{tankID}`
7. Retrieving all sensors for a specific tank
   - `/tanks/{tankID}/tank-sensors`
8. Retrieving all pumps connected to a specific tank of the given id
   - `/tanks/{tankID}/pumps`
9. Changing/updating tank name
   - `/tanks/{tankID}/name`
10. Obtaining meta infor for a given tank
    This endpoint takes Post and Get method
    - `/tanks/{tankID}/profile`
11. Retrieving water levels from a tank
    NB: Kind should be set to _WaterLevel_
    - `/tanks/{tankID}/tank-sensors/waterlevel` -> Water Level Sensor information
    - `/tanks/{tankID}/tank-sensors/waterlevel/value` -> Current value for the water level sensor
    - `/tanks/{tankID}/tank-sensors/waterlevel/values` -> List of water level values
12. Retrieving the pump state from a given device
    NB: Kind should be set to _Motor_
    - `/tanks/{tankID}/pumps/state` -> Shows the recent value
    - `/tanks/{tankID}/pumps/states` -> Includes historical values
13. Perform an actuation
    - One can send a _1_ or a _0_
    - `/tanks/{tankID}/pumps/state`
