package sequencediagram

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	titlePattern       = regexp.MustCompile("^title (.+)$")
	participantPattern = regexp.MustCompile("^participant (.+)$")
	messagePattern     = regexp.MustCompile("^.+->.+:.+$")
	notePattern        = regexp.MustCompile("^note (right|left) of (.+):(.+)$")
)

var arrowRegex = regexp.MustCompile("--?>>?")

func ParseFromText(s string) (*Diagram, error) {
	lines := strings.Split(s, "\n")
	sd := &Diagram{}
	for i, line := range lines {
		switch {
		case titlePattern.MatchString(line):
			title := titlePattern.FindStringSubmatch(line)[1]
			sd.messages = append(sd.messages, Title{simpleMessage{title}})
		case participantPattern.MatchString(line):
			node := sd.getOrCreateNode(participantPattern.FindStringSubmatch(line)[1])
			sd.messages = append(sd.messages, Participant{node, noMessage{}})
		case messagePattern.MatchString(line):
			arrow := arrowRegex.FindString(line)
			message := regexp.MustCompile("^(.+)" + arrow + "(.+):(.+)$").FindStringSubmatch(line)[1:]
			from := sd.getOrCreateNode(message[0])
			to := sd.getOrCreateNode(message[1])
			msg := message[2]
			sd.messages = append(sd.messages, createMessage(from, to, arrow, msg))
		case notePattern.MatchString(line):
			note := notePattern.FindStringSubmatch(line)[1:]
			node := sd.getOrCreateNode(note[1])
			sd.messages = append(sd.messages, createNote(node, note[0], note[2]))
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

// create a Note message
func createNote(node *Node, leftOrRight, msg string) Note {
	side := Left
	if leftOrRight == "right" {
		side = Right
	}
	return Note{node, side, simpleMessage{msg}}
}
