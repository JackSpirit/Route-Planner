package main

func FindIsochroneBoundary(startNodeID int64, maxTimeSeconds float64) [][][]float64 {
	gScore := make(map[int64]float64)
	for nodeID := range networkGraph {
		gScore[nodeID] = 1.7976931348623157e308
	}
	gScore[startNodeID] = 0.0

	pq := NewMinHeapWrapper()
	pq.PushOrUpdate(startNodeID, 0.0)

	var reachableSegments [][][]float64

	for pq.Len() > 0 {
		currentNodeID := pq.Pop()
		currentNode, nodeExists := nodeStorage[currentNodeID]
		if !nodeExists {
			continue
		}

		for _, edge := range networkGraph[currentNodeID] {
			neighborID := edge.ToNodeID
			neighborNode, neighborExists := nodeStorage[neighborID]
			if !neighborExists {
				continue
			}

			weight := CalculateEdgeWeight(edge, "time")
			tentativeGScore := gScore[currentNodeID] + weight

			if tentativeGScore <= maxTimeSeconds {
				if tentativeGScore < gScore[neighborID] {
					gScore[neighborID] = tentativeGScore
					pq.PushOrUpdate(neighborID, tentativeGScore)
				}

				segment := [][]float64{
					{currentNode.Lat, currentNode.Lon},
					{neighborNode.Lat, neighborNode.Lon},
				}
				reachableSegments = append(reachableSegments, segment)
			}
		}
	}

	return reachableSegments
}