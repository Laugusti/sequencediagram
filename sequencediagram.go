// package sequencediagram contains functions to parse a sequence diagram from text
package sequencediagram

import (
	"fmt"
	"sort"
)

type Node struct {
	Name  string
	Order int
}

type Diagram struct {
	messages []Message
	nodes    map[string]*Node
}

func (sd *Diagram) getOrCreateNode(name string) *Node {
	if _, ok := sd.nodes[name]; !ok {
		if sd.nodes == nil {
			sd.nodes = make(map[string]*Node)
		}
		sd.nodes[name] = &Node{name, len(sd.nodes)}
	}
	return sd.nodes[name]
}

func (sd *Diagram) Messages() []Message {
	return sd.messages
}

// GetOrderedNodes returns an ordered Node slice for the sequence diagram
func (sd *Diagram) GetOrderedNodes() []*Node {
	var nodes []*Node
	for _, node := range sd.nodes {
		nodes = append(nodes, node)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Order < nodes[j].Order })
	return nodes
}

// String returns the sequence diagram messages
func (sd *Diagram) String() string {
	var s string
	for _, message := range sd.messages {
		s += fmt.Sprintf("%s\n", message)
	}
	return s
}
