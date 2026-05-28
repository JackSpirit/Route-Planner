package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/route", handleRoute)

	fmt.Println("Web server starting on http://localhost:8080..")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	startLatStr := r.URL.Query().Get("start_lat")
	startLonStr := r.URL.Query().Get("start_lon")
	endLatStr := r.URL.Query().Get("end_lat")
	endLonStr := r.URL.Query().Get("end_lon")

	if startLatStr == "" || startLonStr == "" || endLatStr == "" || endLonStr == "" {
		http.Error(w, `{"error": "Missing coordinates parameters (start_lat, start_lon, end_lat, end_lon required)"}`, http.StatusBadRequest)
		return
	}
    
	startLat, err1 := strconv.ParseFloat(startLatStr, 64)
	startLon, err2 := strconv.ParseFloat(startLonStr, 64)
	endLat, err3 := strconv.ParseFloat(endLatStr, 64)
	endLon, err4 := strconv.ParseFloat(endLonStr, 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, `{"error": "Invalid floating point coordinates"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("Valid coordinates received: Start(%f, %f) -> End(%f, %f)\n", startLat, startLon, endLat, endLon)
    

	startNodeID := findClosestNode(startLat, startLon)
	endNodeID := findClosestNode(endLat, endLon)

	if startNodeID == -1 || endNodeID == -1 {
		http.Error(w, `{"error": "Could not map coordinates to network nodes"}`, http.StatusInternalServerError)
		return
	}

	fmt.Printf("Successfully mapped to graph nodes: Start Node ID %d -> End Node ID %d\n", startNodeID, endNodeID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status": "success", "message": "Mapped coordinates to nodes successfully!", "start_node": %d, "end_node": %d}`, startNodeID, endNodeID)))
}

func findClosestNode(lat, lon float64) int64 {
	var closestNodeID int64 = -1
	minDist := 999999.9

	for nodeID := range networkGraph {
		node, exists := nodeStorage[nodeID]
		if !exists {
			continue
		}

		d := CalculateDistance(lat, lon, node.Lat, node.Lon)
		if d < minDist {
			minDist = d
			closestNodeID = nodeID
		}
	}

	return closestNodeID
}