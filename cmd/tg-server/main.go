package main

import (
	"fmt"
	"github.com/yourbasic/graph"
)

func main() {
	g := graph.New(6)

	path, dist := graph.ShortestPath(g, 0, 5)
	fmt.Println("path:", path, "length:", dist)
}
