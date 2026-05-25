package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
    "github.com/qedus/osmpbf"
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

	var nodeCount int64
	var wayCount int64

	fmt.Println("Streaming binary blocks")

	for {
		v, err := decoder.Decode(); 
		if err == io.EOF {
			break
		} else if err!= nil {
			fmt.Printf("Error during streaming: %v\n", err)
			return
		} else {
			switch v.(type) {
			case *osmpbf.Node:
				nodeCount++
			case *osmpbf.Way:
			    wayCount++
			}
		}
	}

	fmt.Printf("\n Success! Binary Scanning Complete")
	fmt.Printf("Total Nodes Processed: %d\n", nodeCount)
	fmt.Printf("Total Ways Processed: %d\n", wayCount)
}