package main

import (
	"math"
)

type Edge struct {
	ToNodeID int64
	WayID int64
	Distance float64
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

func BuildNetworkGraph() {
	for _, way := range wayStorage {

		for i:=0; i<len(way.NodeIDs)-1; i++ {
			nodeA := way.NodeIDs[i]
			nodeB := way.NodeIDs[i+1]

			nA, foundA := nodeStorage[nodeA]
			nB, foundB := nodeStorage[nodeB]

			if foundA && foundB {
				dist := CalculateDistance(nA.Lat, nA.Lon, nB.Lat, nB.Lon)

				networkGraph[nodeA] = append(networkGraph[nodeA], Edge{
				ToNodeID: nodeB,
				WayID: way.ID,
				Distance: dist,
			    })

			    networkGraph[nodeB] = append(networkGraph[nodeB], Edge{
				ToNodeID: nodeA,
				WayID: way.ID,
				Distance: dist,
			    })
			}
		}
	}
}