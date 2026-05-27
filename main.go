package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	fmt.Println("Parsing OSM PBF file...")
	ParseOSMFile("delaware-260524.osm.pbf")

	fmt.Println("Building network graph...")
	BuildNetworkGraph()

	fmt.Printf("Initialization complete in %v\n", time.Since(start))

	var startNodeID, endNodeID int64
	for _, way := range wayStorage {
		if len(way.NodeIDs) >= 2 {
			startNodeID = way.NodeIDs[0]
			endNodeID = way.NodeIDs[len(way.NodeIDs)-1]
			fmt.Printf("Selected test road: %s (Way ID: %d)\n", way.Name, way.ID)
			break
		}
	}

	if startNodeID == 0 || endNodeID == 0 {
		fmt.Println("Could not find any valid roads with connected nodes to test.")
		return
	}

	fmt.Printf("Finding shortest path from Node %d to Node %d...\n", startNodeID, endNodeID)
	pathRoute := time.Now()
	path := FindShortestPath(startNodeID, endNodeID)
	fmt.Printf("Pathfinding query finished in %v\n", time.Since(pathRoute))

	if path == nil {
		fmt.Println("No path found between the selected nodes.")
		return
	}

	fmt.Printf("Route found! Path contains %d nodes.\n", len(path))
}



