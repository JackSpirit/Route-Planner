package main

import (
	"container/heap"
)

type MinHeapWrapper struct {
	pq    PriorityQueue
	items map[int64]*PathItem
}

func NewMinHeapWrapper() *MinHeapWrapper {
	w := &MinHeapWrapper{
		pq:    make(PriorityQueue, 0),
		items: make(map[int64]*PathItem),
	}
	heap.Init(&w.pq)
	return w
}

func (w *MinHeapWrapper) PushOrUpdate(nodeID int64, distance float64) {
	if item, exists := w.items[nodeID]; exists {
		if distance < item.Distance {
			item.Distance = distance
			heap.Fix(&w.pq, item.index)
		}
	} else {
		item := &PathItem{
			NodeID:   nodeID,
			Distance: distance,
		}
		w.items[nodeID] = item
		heap.Push(&w.pq, item)
	}
}

func (w *MinHeapWrapper) Pop() int64 {
	if len(w.pq) == 0 {
		return -1
	}
	item := heap.Pop(&w.pq).(*PathItem)
	delete(w.items, item.NodeID)
	return item.NodeID
}

func (w *MinHeapWrapper) Len() int {
	return len(w.pq)
}

func FindShortestPath(startNodeID, endNodeID int64, mode string) []int64 {
	gScore := make(map[int64]float64)
	parentMap := make(map[int64]int64)

	for nodeID := range networkGraph {
		gScore[nodeID] = 1.7976931348623157e308
	}
	gScore[startNodeID] = 0.0

	pq := NewMinHeapWrapper()
	pq.PushOrUpdate(startNodeID, 0.0)

	for pq.Len() > 0 {
		currentNodeID := pq.Pop()

		if currentNodeID == endNodeID {
			break
		}

		for _, edge := range networkGraph[currentNodeID] {
			neighborID := edge.ToNodeID
			weight := CalculateEdgeWeight(edge, mode)
			tentativeGScore := gScore[currentNodeID] + weight

			if tentativeGScore < gScore[neighborID] {
				gScore[neighborID] = tentativeGScore
				parentMap[neighborID] = currentNodeID
				pq.PushOrUpdate(neighborID, tentativeGScore)
			}
		}
	}

	var path []int64
	curr := endNodeID
	if _, exists := parentMap[curr]; !exists && curr != startNodeID {
		return nil
	}
	for curr != startNodeID {
		path = append([]int64{curr}, path...)
		curr = parentMap[curr]
	}
	path = append([]int64{startNodeID}, path...)
	return path
}