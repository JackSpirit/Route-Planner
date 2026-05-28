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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success", "message": "Coordinates parsed and validated perfectly. Ready for node lookup!"}`))
}