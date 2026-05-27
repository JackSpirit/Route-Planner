package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	fmt.Println("Parsing OSM PBf file...")
	ParseOSMFile("delaware-260524.osm.pbf")

	fmt.Println("Building network graph...")
	BuildNetworkGraph()

	fmt.Printf("Initialization complete in %v\n", time.Since(start))
}





