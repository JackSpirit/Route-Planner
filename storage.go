package main

type Node struct {
	ID int64
	Lat float64
	Lon float64
}

type Way struct {
	ID int64
	Name string
	NodeIDs []int64
}

var nodeStorage = make(map[int64]Node)
var wayStorage = make(map[int64]Way)

func GetWayCoordinates(wayID int64) [][]float64 {
	way, exists := wayStorage[wayID]
	if !exists {
		return nil
	}

	var coordinates [][]float64
	for _, nodeID := range way.NodeIDs {
		if node, found := nodeStorage[nodeID]; found {
			coordinates = append(coordinates, []float64{node.Lat, node.Lon})
		}
	}
	return coordinates
}