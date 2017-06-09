package sequencediagram

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

func (sd *SequenceDiagram) AsTextDiagram() string {
	var s string
	// Add messages
	// Close system boxes
	drawLifeline := false
	offsets := calcOffsets(sd)
	s += headersAsTextDiagram(sd, offsets) + "\n"
	for i, of := range offsets {
		if i == 0 {
			s += fmt.Sprintf("%s│", strings.Repeat(" ", of.getMiddle()))
		} else {
			s += fmt.Sprintf("%s│", strings.Repeat(" ", of.getMiddle()-offsets[i-1].getMiddle()-1))
		}
	}
	s += "\n"
	for _, message := range sd.Messages {
		if message.isSelfLoop() {
			node := message.From
			for _, line := range strings.Split(selfLoop(message.Msg), "\n") {
				loop := fmt.Sprintf("%s %s", strings.Repeat(" ", offsets[node.Order].getMiddle()), line)
				if drawLifeline {
					loop = fillInLifeline(loop, message, offsets)
				}
				drawLifeline = !drawLifeline
				s += loop + "\n"
			}
		} else {
			if message.From.Order < message.To.Order {
				for i, line := range strings.Split(messageBox(message.Msg), "\n") {
					if i == 1 {
						arrow := strings.Repeat("─", offsets[message.To.Order].getMiddle()-offsets[message.From.Order].getMiddle()-utf8.RuneCountInString(line)-3)
						line = fmt.Sprintf("%s──%s%s▶", strings.Repeat(" ", offsets[message.From.Order].getMiddle()), line, arrow)
					} else {
						line = fmt.Sprintf("%s  %s", strings.Repeat(" ", offsets[message.From.Order].getMiddle()), line)
					}
					if drawLifeline {
						line = fillInLifeline(line, message, offsets)
					}
					drawLifeline = !drawLifeline
					s += line + "\n"
				}
			} else {
				for i, line := range strings.Split(messageBox(message.Msg), "\n") {
					if i == 1 {
						arrow := strings.Repeat("─", offsets[message.From.Order].getMiddle()-offsets[message.To.Order].getMiddle()-utf8.RuneCountInString(line)-4)
						line = fmt.Sprintf("%s◀%s%s──", strings.Repeat(" ", offsets[message.To.Order].getMiddle()+1), arrow, line)
					} else {
						pad := strings.Repeat(" ", offsets[message.From.Order].getMiddle()-offsets[message.To.Order].getMiddle()-utf8.RuneCountInString(line)-1)
						line = fmt.Sprintf("%s  %s", pad, line)
					}
					if drawLifeline {
						line = fillInLifeline(line, message, offsets)
					}
					drawLifeline = !drawLifeline
					s += line + "\n"
				}
			}
		}
	}
	return s
}

func fillInLifeline(s string, m Message, offsets []offset) string {
	var startNode, endNode int
	if m.isSelfLoop() {
		startNode = m.From.Order
		endNode = startNode + 1
		if endNode >= len(offsets) {
			endNode = len(offsets) - 1
		}
	} else {
		startNode = m.From.Order
		endNode = m.To.Order
		if startNode > endNode {
			startNode, endNode = endNode, startNode
		}
	}
	startRange, endRange := offsets[startNode].getMiddle(), offsets[endNode].getMiddle()
	for _, o := range offsets {
		index := o.getMiddle()
		if index <= startRange || index >= endRange {
			if index >= utf8.RuneCountInString(s) {
				s += strings.Repeat(" ", index-utf8.RuneCountInString(s)) + "│"
			} else {
				s = replaceAtIndex(s, index, '│')
			}
		}
	}
	return s
}

func replaceAtIndex(s string, index int, r rune) string {
	var runeSize int
	var runeCount int
	for i := 0; i < len(s); i += runeSize {
		_, runeSize = utf8.DecodeRuneInString(s[i:])
		if runeCount == index {
			s = fmt.Sprintf("%s%s%s", s[:i], string(r), s[i+runeSize:])
			break
		}
		runeCount++
	}
	return s
}

func headersAsTextDiagram(sd *SequenceDiagram, offsets []offset) string {
	nodes := sd.getOrderedNodes()
	height := headerBoxHeight(nodes)
	headers := make([]string, height+2)
	for i, node := range nodes {
		var pad string
		if i != 0 {
			pad = strings.Repeat(" ", offsets[i].begin-offsets[i-1].end-1)
		}
		box := boxString(node.Name, height)
		for j, line := range strings.Split(box, "\n") {
			headers[j] += pad + line
		}
	}
	return strings.Join(headers, "\n")
}

func headerBoxHeight(nodes []*Node) int {
	var max int
	for _, node := range nodes {
		lines := len(strings.Split(node.Name, "\\n"))
		if lines > max {
			max = lines
		}
	}
	return max
}

type offset struct {
	begin int
	end   int
}

func (o offset) getMiddle() int {
	return o.begin + (o.end-o.begin+1)/2
}

func calcOffsets(sd *SequenceDiagram) []offset {
	offsets := make([]offset, len(sd.nodes))
	// calc minimum offsets between nodes
	for i, node := range sd.getOrderedNodes() {
		var begin int
		if i > 0 {
			begin = offsets[i-1].end + 2 + 1
		}
		offsets[i] = offset{begin, begin + len(node.Name) + 4 - 1}
	}

	// adjust offsets based on message
	for _, message := range sd.Messages {
		shiftStart := message.To.Order
		if message.isSelfLoop() {
			shiftStart += 1
		} else if message.From.Order > message.To.Order {
			shiftStart = message.From.Order
		}
		if shiftStart >= len(sd.nodes) {
			continue
		}

		var length int
		for _, m := range message.splitMessage() {
			if len(m) > length {
				length = len(m)
			}
		}

		var shift int
		if message.isSelfLoop() {
			length += 7
			offset1 := offsets[message.From.Order].getMiddle()
			offset2 := offsets[shiftStart].getMiddle()
			diff := offset2 - offset1 - 1
			if length > diff {
				shift = length - diff
			}
		} else {
			length += 8
			offset1 := offsets[message.From.Order].getMiddle()
			offset2 := offsets[message.To.Order].getMiddle()
			diff := int(math.Abs(float64(offset2-offset1))) - 1
			if length > diff {
				shift = length - diff
			}
		}

		for i := shiftStart; i < len(offsets); i++ {
			offsets[i].begin += shift
			offsets[i].end += shift
		}
	}

	return offsets
}

func boxString(s string, padToHeight int) string {
	lines := strings.Split(s, "\\n")
	var maxLength int
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}
	var content string
	for _, line := range lines {
		content += "│" + symmetricPadToLength(line, ' ', maxLength+2) + "│\n"
	}
	middle := strings.Repeat("─", maxLength+2)
	for i := len(lines); i < padToHeight; i++ {
		content += "│" + strings.Repeat("─", maxLength+2) + "│\n"
	}

	box := fmt.Sprintf("┌%s┐\n%s└%s┘", middle, content, middle)
	return box
}

func selfLoop(s string) string {
	loop := "────┐\n"
	for _, line := range strings.Split(s, "\\n") {
		loop += "    │" + line + "\n"
	}
	loop += "◀───┘"
	return loop
}

func messageBox(s string) string {
	box := strings.Split(boxString(s, 0), "\n")
	for i, line := range box {
		if i == 1 {
			box[i] = replaceAtIndex(replaceAtIndex(line, utf8.RuneCountInString(line)-1, '├'), 0, '┤')
		}
	}
	return strings.Join(box, "\n")
}

func symmetricPadToLength(s string, r rune, n int) string {
	if len(s) >= n {
		return s
	}
	padLeft := strings.Repeat(string(r), (n-len(s))/2)
	padRight := strings.Repeat(string(r), (n-len(s)+1)/2)
	return padLeft + s + padRight
}
