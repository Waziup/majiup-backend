package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ApiServe(r *httprouter.Router) {

	// Enable CORS middleware for all endpoints
	r.OPTIONS("/*path", handleOptions)

	// Endpoint to get tanks under majiup
	r.GET("/tanks", handleCORS(TankHandler))

	// Return devices using a specific ID
	r.GET("/tanks/:tankID", handleCORS(GetTankByIDHandler))

	// Endpoint to get all sensors for a specific tank
	r.GET("/tanks/:tankID/tank-sensors", handleCORS(TankSensorHandler))

	r.GET("/tanks/:tankID/tank-info", handleCORS(GetSensorHistoryHandler))

	/*--------------------------------TANK META ENDPOINTS-------------------------------*/

	// r.GET("/tanks/:tankID/location", TankLocationHandler)
	r.POST("/tanks/:tankID/location", handleCORS(TankLocationPostHandler))
	// r.GET("/tanks/:tankID/settings", TankSettingsHandler)
	// r.POST("/tanks/:tankID/settings", TankSettingsPostHandler)

	// r.GET("/tanks/notifications", TanksNotifications)
	// r.GET("/tanks/:tankID/notifications", TankIdNotificationsHandler)
	// r.POST("/tanks/:tankID/notifications", TankNotificationsPostHandler)
	// r.DELETE("/tanks/:tankID/notification", TankNotificationDeleteHandler)

	/*-----------------------------WATER LEVEL SENSOR ENDPOINTS--------------------------------*/

	// Endpoint to get the water level sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/waterlevel", handleCORS(WaterLevelSensorHandler))

	// Endpoint to get the water level value
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/value", handleCORS(GetWaterLevelValueHandler))

	// Endpoint to get the water level history values
	r.GET("/tanks/:tankID/tank-sensors/waterlevel/values", handleCORS(GetWaterLevelHistoryHandler))

	/*-----------------------------WATER TEMPERATURE SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water temperature sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature", handleCORS(WaterTemperatureSensorHandler))

	// Endpoint to get the water temperature value from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature/value", handleCORS(GetWaterTemperatureValueHandler))

	// Endpoint to get the water temperature history values data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-temperature/values", handleCORS(GetWaterTemperatureHistoryHandler))

	/*-----------------------------WATER QUALITY SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water quality sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality", handleCORS(WaterQualitySensorHandler))

	// Endpoint to get the water quality sensor data from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality/value", handleCORS(GetWaterQualityValueHandler))

	// Endpoint to get the water quality history values from a specific tank
	r.GET("/tanks/:tankID/tank-sensors/water-quality/values", handleCORS(GetWaterQualityHistoryHandler))

	/*---------------------------------PUMP ENDPOINTS-------------------------------------------*/
	// Endpoint to get the pumps available in the tank (1)
	r.GET("/tanks/:tankID/pumps", handleCORS(TankPumpHandler))

	// Endpoint to get the pump state of the selected tankID
	r.GET("/tanks/:tankID/pumps/state", handleCORS(TankStateHandler))

	// Endpoint to get the pump states of the selected tankID, the length of the array will give the actuation
	r.GET("/tanks/:tankID/pumps/states", handleCORS(TankStateHistoryHandler))

	// Endpoint to get
	r.POST("/tanks/:tankID/pumps/state", handleCORS(TankStatePostHandler))
}

// handleOptions handles preflight OPTIONS requests
func handleOptions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

// handleCORS wraps a handler function with CORS headers
func handleCORS(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h(w, r, ps)
	}
}
