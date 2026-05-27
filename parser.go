package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
)

func ParseOSMFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open OSM PBF file: %v", err)
	}
	defer file.Close()

	decoder := osmpbf.NewDecoder(file)
	decoder.SetBufferSize(osmpbf.MaxBlobSize)
	
	err = decoder.Start(runtime.NumCPU())
	if err != nil {
		log.Fatalf("Failed to start OSMPBF decoder: %v", err)
	}

	var nc, wc int64

	for {
		v, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error parsing stream: %v", err)
		}

		switch element := v.(type) {
		case *osmpbf.Node:
			nodeStorage[element.ID] = Node{
				ID:  element.ID,
				Lat: element.Lat,
				Lon: element.Lon,
			}
			nc++

		case *osmpbf.Way:
			if hasRoutingTags(element.Tags) {
				wayStorage[element.ID] = Way{
					ID:      element.ID,
					Name:    element.Tags["name"],
					NodeIDs: element.NodeIDs,
				}
				wc++
			}
		}
	}

	log.Printf("Successfully parsed map features. Processed Nodes: %d, Processed Ways: %d\n", nc, wc)
}

func hasRoutingTags(tags map[string]string) bool {
	if _, highwayExists := tags["highway"]; highwayExists {
		return true
	}
	return false
}