package main

import (
	"fmt"
	"math"
)

type Edge struct {
	ToNodeID int64
	WayID    int64
	Distance float64
	SpeedMPH float64
}

var networkGraph = make(map[int64][]Edge)

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000.0

	radLat1 := lat1 * math.Pi/180.0
	radLon1 := lon1 * math.Pi/180.0
	radLat2 := lat2 * math.Pi/180.0
	radLon2 := lon2 * math.Pi/180.0

	diffLat := radLat2-radLat1
	diffLon := radLon2-radLon1

	a := math.Sin(diffLat/2)*math.Sin(diffLat/2)+ math.Cos(radLat1)*math.Cos(radLat2)*math.Sin(diffLon/2)*math.Sin(diffLon/2)

	c := 2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))

	return earthRadius * c
}

func GetSpeedLimit(tags map[string]string) float64 {
	if maxspeedStr, exists := tags["maxspeed"]; exists {
		var speed float64
		_, err := fmt.Sscanf(maxspeedStr, "%f", &speed)
		if err == nil && speed > 0 {
			return speed
		}
	}
	switch tags["highway"] {
	case "motorway":
		return 65.0
	case "trunk", "primary":
		return 55.0
	case "secondary", "tertiary":
		return 45.0
	case "residential":
		return 25.0
	default:
		return 30.0
	}
}

func CalculateEdgeWeight(edge Edge, mode string) float64 {
	if mode == "distance" {
		return edge.Distance
	}
	metersPerSecond := (edge.SpeedMPH * 1609.34) / 3600.0
	if metersPerSecond <= 0 {
		metersPerSecond = 13.4
	}
	return edge.Distance / metersPerSecond
}

func BuildNetworkGraph() {
	for _, way := range wayStorage {
		var speed float64 = 30.0

		for i := 0; i < len(way.NodeIDs)-1; i++ {
			nodeA := way.NodeIDs[i]
			nodeB := way.NodeIDs[i+1]

			nA, foundA := nodeStorage[nodeA]
			nB, foundB := nodeStorage[nodeB]

			if foundA && foundB {
				dist := CalculateDistance(nA.Lat, nA.Lon, nB.Lat, nB.Lon)

				networkGraph[nodeA] = append(networkGraph[nodeA], Edge{
					ToNodeID: nodeB,
					WayID:    way.ID,
					Distance: dist,
					SpeedMPH: speed,
				})

				networkGraph[nodeB] = append(networkGraph[nodeB], Edge{
					ToNodeID: nodeA,
					WayID:    way.ID,
					Distance: dist,
					SpeedMPH: speed,
				})
			}
		}
	}
}