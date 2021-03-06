// Package graph provides an implement of cfn.Graph
// and can be used to graph dependencies between elements
// of a CloudFormation template
package graph

import (
	"fmt"
	"sort"
	"strings"
)

type Graph struct {
	nodes map[interface{}]map[interface{}]bool
	order []interface{}
}

func (g *Graph) String() string {
	out := strings.Builder{}

	for _, left := range g.order {
		if len(g.nodes[left]) > 0 {
			out.WriteString(fmt.Sprintf("%s:\n", left))
			for right, _ := range g.nodes[left] {
				out.WriteString(fmt.Sprintf("- %s\n", right))
			}
		}
	}

	return out.String()
}

// New returns a new Graph
func New() Graph {
	return Graph{
		nodes: make(map[interface{}]map[interface{}]bool),
	}
}

func (g *Graph) add(item interface{}) {
	if _, ok := g.nodes[item]; !ok {
		g.nodes[item] = make(map[interface{}]bool)
		g.order = append(g.order, item)
	}
}

// Add creates a link between two nodes in the graph
func (g *Graph) Add(item interface{}, links ...interface{}) {
	g.add(item)

	for _, to := range links {
		g.add(to)
		g.nodes[item][to] = true
	}
}

func (g Graph) depth(item interface{}) int {
	seen := map[interface{}]bool{
		item: true,
	}

	count := 0

	var dive func(interface{})

	dive = func(from interface{}) {
		for to, _ := range g.nodes[from] {
			if !seen[to] {
				seen[to] = true
				count++

				dive(to)
			}
		}
	}

	dive(item)

	return count
}

// Nodes returns all nodes of the graph, in order of their dependencies.
// Nodes with the fewest dependencies are at the beginning of the slice.
func (g Graph) Nodes() []interface{} {
	sort.Slice(g.order, func(i, j int) bool {
		a, b := g.order[i], g.order[j]

		diff := g.depth(a) - g.depth(b)

		if diff == 0 {
			return fmt.Sprint(a) < fmt.Sprint(b)
		}

		return diff < 0
	})

	return g.order
}

// Get returns all nodes that are connected to the item that you pass in.
func (g Graph) Get(item interface{}) []interface{} {
	links := make([]interface{}, 0)
	for to, _ := range g.nodes[item] {
		links = append(links, to)
	}

	sort.Slice(links, func(i, j int) bool {
		return fmt.Sprint(links[i]) < fmt.Sprint(links[j])
	})

	return links
}

// GetReverse returns all nodes that connect to the item that you pass in.
func (g Graph) GetReverse(item interface{}) []interface{} {
	links := make([]interface{}, 0)
	for from, deps := range g.nodes {
		if _, ok := deps[item]; ok {
			links = append(links, from)
		}
	}

	sort.Slice(links, func(i, j int) bool {
		return fmt.Sprint(links[i]) < fmt.Sprint(links[j])
	})

	return links
}
