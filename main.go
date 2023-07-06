package main

import (
	"fmt"
	"net/http"

	"github.com/JosephMusya/majiup-backend/api"
	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()

	/*----------------------------------TANK ENDPOINTS-------------------------------*/

	// Endpoint to get tanks under majiup
	router.GET("/tanks", api.TankHandler)

	// Return devices using a specific ID
	router.GET("/tanks/:tankID", api.GetTankByIDHandler)

	// Endpoint to get all sensors for a specific tank
	router.GET("/tanks/:tankID/tank-sensors", api.TankSensorHandler)

	/*--------------------------------TANK META ENDPOINTS-------------------------------*/

	/*-----------------------------WATER LEVEL SENSOR ENDPOINTS--------------------------------*/

	// Endpoint to get the water level sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/waterlevel", api.WaterLevelSensorHandler)

	// Endpoint to get the water level value
	router.GET("/tanks/:tankID/tank-sensors/waterlevel/value", api.GetWaterLevelValueHandler)

	// Endpoint to get the water level history values
	router.GET("/tanks/:tankID/tank-sensors/waterlevel/values", api.GetWaterLevelHistoryHandler)

	/*-----------------------------WATER TEMPERATURE SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water temperature sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature", api.WaterTemperatureSensorHandler)

	// Endpoint to get the water temperature value from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature/value", api.GetWaterTemperatureValueHandler)

	// Endpoint to get the water temperature history values data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-temperature/values", api.GetWaterTemperatureHistoryHandler)

	/*-----------------------------WATER QUALITY SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water quality sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality", api.WaterQualitySensorHandler)

	// Endpoint to get the water quality sensor data from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality/value", api.GetWaterQualityValueHandler)

	// Endpoint to get the water quality history values from a specific tank
	router.GET("/tanks/:tankID/tank-sensors/water-quality/values", api.GetWaterQualityHistoryHandler)

	/*---------------------------------PUMP ENDPOINTS-------------------------------------------*/
	// Endpoint to get the pumps available in the tank (1)
	router.GET("/tanks/:tankID/pumps", api.TankPumpHandler)

	// Endpoint to get the pump state of the selected tankID
	router.GET("/tanks/:tankID/pumps/state", api.TankStateHandler)

	// Endpoint to get the pump states of the selected tankID, the length of the array will give the actuation
	router.GET("/tanks/:tankID/pumps/states", api.TankStateHistoryHandler)

	// Endpoint to get
	router.POST("/tanks/:tankID/pumps/state", api.TankStatePostHandler)

	fmt.Println("Majiup server running at PORT 8080")
	http.ListenAndServe(":8080", router)
}
