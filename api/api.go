package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ApiServe(r *mux.Router) {
	// Enable CORS middleware for all endpoints
	r.HandleFunc("/{path:.*}", handleOptions).Methods("OPTIONS")

	// Gateway profile
	r.HandleFunc("/gateway-profile", handleCORS(getGatewayProfile)).Methods("GET")
	r.HandleFunc("/gateway-profile", handleCORS(updateGatewayProfile)).Methods("POST")
	r.HandleFunc("/wifi-status", handleCORS(getWifiStatus)).Methods("GET")	

	// Get analytics
	r.HandleFunc("/tanks/{tankID}/analytics", handleCORS(getAnalytics)).Methods("GET")

	// Internet & services
	// r.HandleFunc("/internet", handleCORS(getWifiStatus)).Methods("GET")

	// Ask majiup copilot
	// r.HandleFunc("/ask-majiup-copilot", handleCORS(AskMajiupCopilot)).Methods("POST")

	// Endpoint to get tanks under majiup
	r.HandleFunc("/tanks", handleCORS(TankHandler)).Methods("GET")

	// Return devices using a specific ID
	r.HandleFunc("/tanks/{tankID}", handleCORS(GetTankByIDHandler)).Methods("GET")

	// Endpoint to get all sensors for a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors", handleCORS(TankSensorHandler)).Methods("GET")

	// Endpoint to get the pumps available in the tank
	r.HandleFunc("/tanks/{tankID}/pumps", handleCORS(TankPumpHandler)).Methods("GET")

	// Endpoint to get sensor history
	r.HandleFunc("/tanks/{tankID}/tank-info", handleCORS(GetSensorHistoryHandler)).Methods("GET")

	// Endpoint to change the name of a device
	r.HandleFunc("/tanks/{tankID}/name", handleCORS(ChangeNameHandler)).Methods("POST")

	// Endpoint to delete a tank
	r.HandleFunc("/tanks/{tankID}", handleCORS(DeleteTank)).Methods("DELETE")

	/*--------------------------------TANK META ENDPOINTS-------------------------------*/

	// GET Meta fields (settings & notifications)
	r.HandleFunc("/tanks/{tankID}/profile", handleCORS(getMetaFields)).Methods("GET")

	// POST Meta fields
	r.HandleFunc("/tanks/{tankID}/profile", handleCORS(postMetaField)).Methods("POST")

	/*-----------------------------WATER LEVEL SENSOR ENDPOINTS--------------------------------*/

	// Endpoint to get the water level sensor data from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/waterlevel", handleCORS(WaterLevelSensorHandler)).Methods("GET")

	// Endpoint to get the water level value
	r.HandleFunc("/tanks/{tankID}/tank-sensors/waterlevel/value", handleCORS(GetWaterLevelValueHandler)).Methods("GET")

	// Endpoint to get the water level history values
	r.HandleFunc("/tanks/{tankID}/tank-sensors/waterlevel/values", handleCORS(GetWaterLevelHistoryHandler)).Methods("GET")

	// Endpoint to change the waterlevel meta field
	r.HandleFunc("/tanks/{tankID}/tank-sensors/waterlevel/alerts", handleCORS(ChangeWaterLevelAlerts)).Methods("POST")

	/*-----------------------------WATER TEMPERATURE SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water temperature sensor data from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-temperature", handleCORS(WaterTemperatureSensorHandler)).Methods("GET")

	// Endpoint to get the water temperature value from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-temperature/value", handleCORS(GetWaterTemperatureValueHandler)).Methods("GET")

	// Endpoint to get the water temperature history values data from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-temperature/values", handleCORS(GetWaterTemperatureHistoryHandler)).Methods("GET")

	// Endpoint to change the water temp level
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-temperature/alerts", handleCORS(ChangeWaterTemperatureAlerts)).Methods("POST")

	/*-----------------------------WATER QUALITY SENSOR ENDPOINTS---------------------------*/

	// Endpoint to get the water quality sensor data from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-quality", handleCORS(WaterQualitySensorHandler)).Methods("GET")

	// Endpoint to get the water quality sensor data from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-quality/value", handleCORS(GetWaterQualityValueHandler)).Methods("GET")

	// Endpoint to get the water quality history values from a specific tank
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-quality/values", handleCORS(GetWaterQualityHistoryHandler)).Methods("GET")

	// Endpoint to change water quality alerts
	r.HandleFunc("/tanks/{tankID}/tank-sensors/water-quality/alerts", handleCORS(ChangeWaterQualityAlerts)).Methods("POST")

	/*---------------------------------PUMP ENDPOINTS-------------------------------------------*/

	// Endpoint to get the pump state of the selected tankID
	r.HandleFunc("/tanks/{tankID}/pumps/state", handleCORS(TankStateHandler)).Methods("GET")

	// Endpoint to get the pump states of the selected tankID, the length of the array will give the actuation
	r.HandleFunc("/tanks/{tankID}/pumps/states", handleCORS(TankStateHistoryHandler)).Methods("GET")

	// Endpoint to post pump state
	r.HandleFunc("/tanks/{tankID}/pumps/state", handleCORS(TankStatePostHandler)).Methods("POST")

	// Handle undefined routes
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the timestamp and the endpoint that was not found
		log.Printf(" [ ERR ][%s] Method not allowed: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})
}

// handleOptions handles preflight OPTIONS requests
func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

// handleCORS wraps a handler function with CORS headers
func handleCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h(w, r)
	}
}
