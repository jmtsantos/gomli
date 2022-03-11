package cmd

import (
	"fmt"
)

func Find(className string) {

	// Search
	minigraph := appGraph.edges[Node{className}]
	fmt.Printf("%v\n", appGraph.edges[Node{className}])

	for _, nodes := range minigraph {
		for _, v := range app[nodes.String()].Methods {
			fmt.Printf("%s methods: %v\n", nodes.String(), v.Name)
		}
	}
}

// Transverse
// appGraph.Traverse(func(n *Node) {
// 	if n.String() == "Lcom/i/i/i/i/k;" {
// 		fmt.Printf("%v %v\n", n, appGraph.edges[*n])
// 	}
// })
