package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RouteResponse struct {
	Path       []Coordinate `json:"path"`
	DistanceKM float64      `json:"distance_km"`
	NodeCount  int          `json:"node_count"`
}

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
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	startLatStr := r.URL.Query().Get("start_lat")
	startLonStr := r.URL.Query().Get("start_lon")
	endLatStr := r.URL.Query().Get("end_lat")
	endLonStr := r.URL.Query().Get("end_lon")
	mode := r.URL.Query().Get("mode")

	if mode != "distance" && mode != "time" {
		mode = "distance"
	}

	if startLatStr == "" || startLonStr == "" || endLatStr == "" || endLonStr == "" {
		http.Error(w, `{"error": "Missing coordinates parameters"}`, http.StatusBadRequest)
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

	startNodeID := findClosestNode(startLat, startLon)
	endNodeID := findClosestNode(endLat, endLon)

	if startNodeID == -1 || endNodeID == -1 {
		http.Error(w, `{"error": "Could not map coordinates to network nodes"}`, http.StatusInternalServerError)
		return
	}

	nodePath := FindShortestPath(startNodeID, endNodeID, mode)
	if nodePath == nil {
		http.Error(w, `{"error": "No continuous path found between these locations"}`, http.StatusNotFound)
		return
	}

	var coordinates []Coordinate
	var totalDistance float64
	var prevNode Node

	for i, nodeID := range nodePath {
		node := nodeStorage[nodeID]
		coordinates = append(coordinates, Coordinate{
			Lat: node.Lat,
			Lon: node.Lon,
		})

		if i > 0 {
			totalDistance += CalculateDistance(prevNode.Lat, prevNode.Lon, node.Lat, node.Lon)
		}
		prevNode = node
	}

	response := RouteResponse{
		Path:       coordinates,
		DistanceKM: totalDistance/1000.0,
		NodeCount:  len(coordinates),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func findClosestNode(lat, lon float64) int64 {
	var closestNodeID int64 = -1
	minDist := 1.7976931348623157e308
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