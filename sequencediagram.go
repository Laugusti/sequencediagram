package sequencediagram

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type Message struct {
	From *Node
	To   *Node
	Msg  string
}

type Node struct {
	Name  string
	Order int
}

type SequenceDiagram struct {
	Messages []Message
	nodes    map[string]*Node
}

func (m Message) splitMessage() []string {
	return strings.Split(m.Msg, "\\n")
}

func (m Message) isSelfLoop() bool {
	return m.From == m.To
}

func (sd *SequenceDiagram) getOrCreateNode(name string) *Node {
	if _, ok := sd.nodes[name]; !ok {
		if sd.nodes == nil {
			sd.nodes = make(map[string]*Node)
		}
		sd.nodes[name] = &Node{name, len(sd.nodes)}
	}
	return sd.nodes[name]
}

type NodeSorter struct {
	nodes []*Node
}

func (s NodeSorter) Len() int           { return len(s.nodes) }
func (s NodeSorter) Swap(i, j int)      { s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i] }
func (s NodeSorter) Less(i, j int) bool { return s.nodes[i].Order < s.nodes[j].Order }

func (sd *SequenceDiagram) getOrderedNodes() []*Node {
	var nodes []*Node
	for _, node := range sd.nodes {
		nodes = append(nodes, node)
	}
	sort.Sort(NodeSorter{nodes})
	return nodes
}

func (sd *SequenceDiagram) String() string {
	var s string
	for _, message := range sd.Messages {
		s += fmt.Sprintf("%s->%s:%s\n", message.From.Name, message.To.Name, message.Msg)
	}
	return s
}

// var titlePattern = regexp.MustCompile("^title .+$")
// var notePattern = regexp.MustCompile("^note (right|left) of .+$")
var messagePattern = regexp.MustCompile("^.+->.+:.+$")

func ParseFromText(s string) (*SequenceDiagram, error) {
	lines := strings.Split(s, "\n")
	sd := &SequenceDiagram{}
	for i, line := range lines {
		switch {
		case messagePattern.MatchString(line):
			arrowIndex := strings.Index(line, "->")
			colonIndex := strings.Index(line, ":")
			from := sd.getOrCreateNode(line[:arrowIndex])
			to := sd.getOrCreateNode(line[arrowIndex+2 : colonIndex])
			msg := line[colonIndex+1:]
			sd.Messages = append(sd.Messages, Message{from, to, msg})
		default:
			return nil, fmt.Errorf("Line %d: Syntax error.", i)
		}
	}
	return sd, nil
}
