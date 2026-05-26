package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
    "github.com/qedus/osmpbf"
)

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

type Edge struct {
	ToNodeID int64
	WayID int64
	Distance float64
}

var (
	nodeStorage = make(map[int64]Node)
	wayStorage = make(map[int64]Way)
	networkGraph = make(map[int64][]Edge)
)

func main() {
	file, err := os.Open("delaware-260524.osm.pbf")
	if err != nil {
		fmt.Printf("Error opening PBF file: %v\n", err)
		return
	}
	defer file.Close()

	decoder := osmpbf.NewDecoder(file)

	decoder.SetBufferSize(osmpbf.MaxBlobSize)
	err = decoder.Start(runtime.NumCPU())
	if err != nil {
		fmt.Printf("Error starting decoder: %v\n", err)
		return
	}

	fmt.Println("Streaming binary blocks")

	for {
		v, err := decoder.Decode(); 
		if err == io.EOF {
			break
		} else if err!= nil {
			fmt.Printf("Error during streaming: %v\n", err)
			return
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				
				nodeStorage[v.ID] = Node{
					ID: v.ID,
					Lat: v.Lat,
					Lon: v.Lon,
				}
			case *osmpbf.Way:
			    streetName := v.Tags["name"]

				wayStorage[v.ID] = Way{
					ID: v.ID,
					Name: streetName,
					NodeIDs: v.NodeIDs,
				}
			}
		}
	}

	fmt.Printf("\n Success! Binary Scanning Complete\n")
	fmt.Printf("Total Nodes cached: %d\n", len(nodeStorage))
	fmt.Printf("Total Ways cached: %d\n", len(wayStorage))
}

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