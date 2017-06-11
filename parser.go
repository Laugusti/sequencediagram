package sequencediagram

import (
	"fmt"
	"regexp"
	"strings"
)

// var titlePattern = regexp.MustCompile("^title .+$")
// var notePattern = regexp.MustCompile("^note (right|left) of .+$")
var messagePattern = regexp.MustCompile("^.+->.+:.+$")

func ParseFromText(s string) (*Diagram, error) {
	lines := strings.Split(s, "\n")
	sd := &Diagram{}
	for i, line := range lines {
		switch {
		case messagePattern.MatchString(line):
			arrowIndex := strings.Index(line, "->")
			colonIndex := strings.Index(line, ":")
			from := sd.getOrCreateNode(line[:arrowIndex])
			to := sd.getOrCreateNode(line[arrowIndex+2 : colonIndex])
			msg := line[colonIndex+1:]
			sd.messages = append(sd.messages, Message{from, to, msg})
		default:
			return nil, fmt.Errorf("Line %d: Syntax error.", i+1)
		}
	}
	return sd, nil
}
