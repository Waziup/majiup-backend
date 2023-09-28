# Majiup

## Steps to run the Majiup
-----  -----

The frontend is served from the backend, therefore, there is no need to download the frontend files.
Note that the application is running on the wazigate ip address, to access the application you must have your gateway powered.
The application can be accessed from https://wazigate.local:8081
****In place of wazigate.local, you can key in the ip address of your gateway.***

### Step 0 - SSH to gateway
```
ssh pi@wazigate.local
```

### Step 1 - Clone the repository
```
git clone https://github.com/JosephMusya/majiup-backend.git
```
### Step 2 - Navigate to the majiup-backend repository
```
cd majiup-backend
```
### Step 3 - Build Majiup image to run on the wazigate
```
sudo docker build --platform linux/arm64  -t majiup .
```
If you are building the image on the gateway itself, run the following command to start the application

### Step 4 - Start the majiup application on the gateway
The application runs on detached mode
```
sudo docker-compose up -d
```
You can check for any messages and troubleshooting approaches by
```
sudo docker logs <majiup-container>
```
If you are on a separed computer, you will need ssh and ftp enabled in your gateway to transfer the build image.

### Step 5 - Save the image in a zip folder
```
sudo docker save -o majiup.tar majiup
```
This saves the image into majiup.tar compressed folder
You can confirm the image folder with
```
ls
```

### Step 6 - Change the write and read permission for the folder
```
sudo chmod 777 majiup.tar
```

### Step 7 - Transfer the folder to gateway
```
ftp <IP_ADDRESS>
```
The IP address is the gateway's ip address
The default username is ***pi*** and password is ***loragateway***

Transfer the file with
```
put majiup.tar
```
### Step 8 - SSH into the raspberry pi
```
ssh pi@<IP_ADDRESS>
```
Load majiup image from the compressed folder
```
sudo docker image load -i majiup.tar
```

### Step 9 - Run the application
This step is similar to step 4
```
sudo docker-compose up -d
```
The application can be accessed from http://wazigate.local:8081

The api is served by http://wazigate.local:8081/api/

## Creating your Majiup device
### Step 1 - Open the wazigate UI
Open the gateway UI dashboard from http://wazigate.local and navigate to the dashbaord
Create a new device and assign unique adresses to the device
Allocate one actuator, the other three sensors will be allocated automatically when the hardware mounted on the tank sends data

### Step 2 -  Setup your hardware
Upload this code to your hardware. Normally, the sensor sends data at inteval of 5 minutes.
Note: Edit the sensor pins defined in the code and the device address as connected to your hardware before uploading.
e.g In the code, the TDS sensor pin is A1, temperature probe connected to pin A2, echo pin and trigger pins are D3 and D4 respectively
When done changing the necessary fields, upload your code and set the device ready to make measurements. Turn the device on when you are done mounting to avoid streaming false value. Even though, this has been catered in the code to reject sending outlier values (false data).

### Step 3 - Set sensors on the gateway
After receiving sensor values on the gateway, edit their fields respectively
Normall,the sensors will have name like temperature sensor 1,....
Sensor with ID of temperature_sensor_0 is water level sensor -> Assign kind to WaterLevel
Sensor with ID of temperature_sensor_1 is water temperature sensor -> Assign kind to WaterThermometer
Sensor with ID of temperature_sensor_2 is water quality sensor sensor -> Assign kind to WaterPollutant

### Step 4 - Set your tank on Majiup application
Head over to majiup application and set the dimensions and capacity of your tank under settings.
