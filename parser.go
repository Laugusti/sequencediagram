package sequencediagram

import (
	"fmt"
	"regexp"
	"strings"
)

// var titlePattern = regexp.MustCompile("^title .+$")
// var notePattern = regexp.MustCompile("^note (right|left) of .+$")
var messagePattern = regexp.MustCompile("^.+->.+:.+$")

var arrowRegex = regexp.MustCompile("--?>>?")

func ParseFromText(s string) (*Diagram, error) {
	lines := strings.Split(s, "\n")
	sd := &Diagram{}
	for i, line := range lines {
		switch {
		case messagePattern.MatchString(line):
			arrow := arrowRegex.FindString(line)
			message := regexp.MustCompile("^(.+)" + arrow + "(.+):(.+)$").FindStringSubmatch(line)[1:]
			from := sd.getOrCreateNode(message[0])
			to := sd.getOrCreateNode(message[1])
			msg := message[2]
			sd.messages = append(sd.messages, createMessage(from, to, arrow, msg))
		default:
			return nil, fmt.Errorf("Line %d: Syntax error.", i+1)
		}
	}
	return sd, nil
}

// creates a self/from/to message
func createMessage(from, to *Node, arrowType, msg string) Message {
	altArrowBody := strings.HasPrefix(arrowType, "--")
	altArrowEnd := strings.HasSuffix(arrowType, ">>")
	switch {
	case from.Order == to.Order:
		return SelfMessage{from, simpleMessage{msg}, uniDirectionalMessage{altArrowBody, altArrowEnd}}
	case from.Order < to.Order:
		return ForwardMessage{from, to, simpleMessage{msg}, uniDirectionalMessage{altArrowBody, altArrowEnd}}
	case from.Order > to.Order:
		return BackwardMessage{from, to, simpleMessage{msg}, uniDirectionalMessage{altArrowBody, altArrowEnd}}
	}
	return nil
}
