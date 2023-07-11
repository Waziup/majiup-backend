package api

import "github.com/julienschmidt/httprouter"

func ApiServe(r *httprouter.Router) {
	// Endpoint to get tanks under majiup
	r.GET("/tanks", TankHandler)

	// Return devices using a specific ID
	r.GET("/tanks/:tankID", GetTankByIDHandler)

	// Endpoint to get all sensors for a specific tank
	r.GET("/tanks/:tankID/tank-sensors", TankSensorHandler)

	r.GET("/tanks/:tankID/tank-info", GetSensorHistoryHandler)

	/*--------------------------------TANK META ENDPOINTS-------------------------------*/

	// r.GET("/tanks/:tankID/location", TankLocationHandler)
	r.POST("/tanks/:tankID/location", TankLocationPostHandler)
	// r.GET("/tanks/:tankID/settings", TankSettingsHandler)
	// r.POST("/tanks/:tankID/settings", TankSettingsPostHandler)

	// r.GET("/tanks/notifications", TanksNotifications)
	// r.GET("/tanks/:tankID/notifications", TankIdNotificationsHandler)
	// r.POST("/tanks/:tankID/notifications", TankNotificationsPostHandler)
	// r.DELETE("/tanks/:tankID/notification", TankNotificationDeleteHandler)

	/*-----------------------------WATER LEVEL SENSOR ENDPOINTS--------------------------------*/

	// Endpoint to get the water level sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/waterlevel", WaterLevelSensorHandler)

	// Endpoint to get the water level value
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/value", GetWaterLevelValueHandler)

	// Endpoint to get the water level history values
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/values", GetWaterLevelHistoryHandler)

	/*-----------------------------WATER TEMPERATURE SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water temperature sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature", WaterTemperatureSensorHandler)

	// Endpoint to get the water temperature value from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature/value", GetWaterTemperatureValueHandler)

	// Endpoint to get the water temperature history values data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature/values", GetWaterTemperatureHistoryHandler)

	/*-----------------------------WATER QUALITY SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water quality sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality", WaterQualitySensorHandler)

	// Endpoint to get the water quality sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality/value", GetWaterQualityValueHandler)

	// Endpoint to get the water quality history values from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality/values", GetWaterQualityHistoryHandler)

	/*---------------------------------PUMP ENDPOINTS-------------------------------------------*/
	// Endpoint to get the pumps available in the tank (1)
	r.GET("/tanks/:tankID/pumps", TankPumpHandler)

	// Endpoint to get the pump state of the selected tankID
	r.GET("/tanks/:tankID/pumps/state", TankStateHandler)

	// Endpoint to get the pump states of the selected tankID, the length of the array will give the actuation
	r.GET("/tanks/:tankID/pumps/states", TankStateHistoryHandler)

	// Endpoint to get
	r.POST("/tanks/:tankID/pumps/state", TankStatePostHandler)
}
