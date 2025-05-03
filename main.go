package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Edge struct {
	Node   string
	Weight int
}

type Graph struct {
	edges map[string][]Edge
	nodes map[string]struct{}
}

func NewGraph() *Graph {
	return &Graph{
		edges: make(map[string][]Edge),
		nodes: make(map[string]struct{}),
	}
}

func (g *Graph) AddEdge(src, dest string, weight int) {
	if weight < 0 {
		panic("negative weight detected")
	}
	g.addNode(src)
	g.addNode(dest)
	g.edges[src] = append(g.edges[src], Edge{Node: dest, Weight: weight})
}

func (g *Graph) addNode(node string) {
	if _, exists := g.nodes[node]; !exists {
		g.nodes[node] = struct{}{}
	}
}

type Item struct {
	node     string
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func (g *Graph) Dijkstra(start string) (distances map[string]int, predecessors map[string]string, err error) {
	if _, exists := g.nodes[start]; !exists {
		return nil, nil, fmt.Errorf("start node '%s' not found", start)
	}

	distances = make(map[string]int)
	predecessors = make(map[string]string)
	for node := range g.nodes {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Item{node: start, priority: 0})

	processed := make(map[string]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item)
		if processed[current.node] {
			continue
		}
		processed[current.node] = true

		for _, edge := range g.edges[current.node] {
			newDist := distances[current.node] + edge.Weight
			if newDist < distances[edge.Node] {
				distances[edge.Node] = newDist
				predecessors[edge.Node] = current.node
				heap.Push(&pq, &Item{
					node:     edge.Node,
					priority: newDist,
				})
			}
		}
	}

	return distances, predecessors, nil
}

func GetPath(predecessors map[string]string, start, target string) ([]string, error) {
	path := []string{}
	current := target

	for {
		path = append([]string{current}, path...)
		if current == start {
			break
		}
		prev, exists := predecessors[current]
		if !exists {
			return nil, fmt.Errorf("no path from %s to %s", start, target)
		}
		current = prev
	}

	if len(path) == 0 || path[0] != start {
		return nil, fmt.Errorf("invalid path")
	}

	return path, nil
}

func main() {
	graph := NewGraph()
	graph.AddEdge("A", "B", 4)
	graph.AddEdge("A", "C", 2)
	graph.AddEdge("B", "C", 5)
	graph.AddEdge("B", "D", 10)
	graph.AddEdge("C", "D", 3)
	graph.AddEdge("C", "E", 1)
	graph.AddEdge("D", "E", 4)

	start := "A"
	distances, predecessors, err := graph.Dijkstra(start)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Shortest distances from", start)
	for node, dist := range distances {
		fmt.Printf("%s: %d\n", node, dist)
	}

	target := "E"
	path, err := GetPath(predecessors, start, target)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Shortest path from %s to %s: %v\n", start, target, path)
}
