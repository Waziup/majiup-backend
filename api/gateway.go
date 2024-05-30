package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Gateway struct {
	Profile		Profile `json:"profile" bson:"profile"`
	// Created 	string	`json:"created" bson:"created"`
	// Id			string 	`json:"id" bson:"id"`
	// Name		string 	`json:"name" bson:"name"`
}

// func getWifiStatus(w http.ResponseWriter, r *http.Request) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", "http://wazigate.local/apps/waziup.wazigate-system/internet", nil)
// 	if err != nil {
// 		log.Println("Error creating HTTP request:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the Content-Type and Accept headers
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Accept", "application/json")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		log.Println("Error obtaining gateway internet:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Check the HTTP response status code
// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := ioutil.ReadAll(resp.Body)
// 		log.Printf("Error: received non-200 status code %d. Response: %s", resp.StatusCode, string(body))
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Read the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Error reading response body:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}


// 	var status bool
// 	err = json.Unmarshal(body, &status)
// 	if err != nil {
// 		log.Println("Error unmarshaling gateway:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Marshal the gateway struct into JSON
// 	response, err := json.Marshal(status)
// 	if err != nil {
// 		log.Println("Error marshaling gateway:", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the Content-Type header to application/json
// 	w.Header().Set("Content-Type", "application/json")

// 	// Log the success message
// 	log.Printf("[%s] Fetched internet status %s: %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

// 	// Write the JSON response to the response writer
// 	w.Write(response)	
// }

func getGatewayProfile(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost/device/meta", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type and Accept headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error obtaining gateway profile:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the HTTP response status code
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error: received non-200 status code %d. Response: %s", resp.StatusCode, string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var gateway Gateway
	err = json.Unmarshal(body, &gateway)
	if err != nil {
		log.Println("Error unmarshaling gateway:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the gateway struct into JSON
	response, err := json.Marshal(gateway)
	if err != nil {
		log.Println("Error marshaling gateway:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Log the success message
	log.Printf("[%s] Fetched gateway %s: %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	// Write the JSON response to the response writer
	w.Write(response)
}

func updateGatewayProfile(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// Send a POST request to localhost/device/meta
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPost, "http://localhost/device/meta", bytes.NewReader(body))
	if err != nil {
		fmt.Println("Error creating POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to text/plain
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Received non-OK response:", resp.Status)
		// Print the response body for further analysis
		responseBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response body:", string(responseBody))
		w.WriteHeader(resp.StatusCode)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write a success response
	response := map[string]string{
		"message": "Gateway profile updated successfully",
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("[%s] Gateway meta field updated successfully: %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

	w.Write(responseBytes)
}