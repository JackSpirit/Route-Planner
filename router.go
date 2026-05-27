package main

import (
	"container/heap"
	"math"
)

func FindShortestPath(startNodeID, endNodeID int64) []int64 {
	if startNodeID == endNodeID {
		return []int64{startNodeID}
	}

	dist := make(map[int64]float64)
	parent := make(map[int64]int64)
	
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	dist[startNodeID] = 0.0
	heap.Push(&pq, &PathItem{
		NodeID:   startNodeID,
		Distance: 0.0,
	})

	var found bool

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PathItem)
		u := current.NodeID

		if u == endNodeID {
			found = true
			break
		}

		if current.Distance > dist[u] {
			continue
		}

		edges := networkGraph[u]
		for _, edge := range edges {
			v := edge.ToNodeID
			weight := edge.Distance

			currentDist, exists := dist[v]
			if !exists {
				currentDist = math.Inf(1)
			}

			if dist[u]+weight < currentDist {
				dist[v] = dist[u] + weight
				parent[v] = u
				heap.Push(&pq, &PathItem{
					NodeID:   v,
					Distance: dist[v],
				})
			}
		}
	}

	if !found {
		return nil
	}

	var path []int64
	curr := endNodeID
	for curr != startNodeID {
		path = append(path, curr)
		curr = parent[curr]
	}
	path = append(path, startNodeID)

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}