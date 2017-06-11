//package textdiagram provides functionality to create a textual representation of a sequence diagram
package textdiagram

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/Laugusti/sequencediagram"
)

type textDiagram struct {
	offsets        []offset
	lifelineToggle bool
	text           string
}

// Decode creates an textual representation a sequence diagram using the
// provided sequence diagram
func Decode(sd *sequencediagram.Diagram) io.Reader {
	td := &textDiagram{}
	td.offsets = calcOffsets(sd)

	nodes := sd.GetOrderedNodes()
	td.addHeaders(nodes)
	td.drawFullLifeline()
	for _, message := range sd.Messages() {
		td.addMessage(message)
	}
	if !td.lifelineToggle {
		td.text += "\n"
	}
	td.drawFullLifeline()
	td.addHeaders(nodes)
	return strings.NewReader(td.text)
}

// addHeaders add the Node slice as text to the ascii diagram
func (td *textDiagram) addHeaders(nodes []*sequencediagram.Node) {
	// get max # of lines in Node Name
	height := headerBoxHeight(nodes)
	// headers will have max + 2 lines (top + #lines + bottom)
	headers := make([]string, height+2)
	for i, node := range nodes {
		var pad string
		if i != 0 {
			// add padding using pre-calculated node offsets
			pad = strings.Repeat(" ", td.offsets[i].begin-td.offsets[i-1].end-1)
		}
		// add each line of box to header slice with padding
		box := boxString(node.Name, height)
		for j, line := range strings.Split(box, "\n") {
			headers[j] += pad + line
		}
	}
	td.text += strings.Join(headers, "\n") + "\n"
}

// addMessage adds the message the as text to the ascii diagram
func (td *textDiagram) addMessage(message sequencediagram.Message) {
	// pad with spaces
	pad := strings.Repeat(" ", td.offsets[message.From.Order].getMiddle()+utf8.RuneCountInString(box_vertical))
	// backward messages should use middle of 'to'
	if message.From.Order > message.To.Order {
		pad = strings.Repeat(" ", td.offsets[message.To.Order].getMiddle()+utf8.RuneCountInString(box_vertical))
	}

	text := td.getMessageAsText(message)
	for _, line := range strings.Split(text, "\n") {
		//for each line, pad and draw life lines
		line = td.fillInLifeline(pad+line, message)
		td.text += line + "\n"
	}
}

// returns the text representation of the message
func (td *textDiagram) getMessageAsText(message sequencediagram.Message) string {
	var text string
	switch {
	case message.IsSelfLoop():
		text = selfLoop(message.Msg)
	case message.From.Order < message.To.Order:
		text = td.toMessageText(message)
	case message.From.Order > message.To.Order:
		text = td.fromMessageText(message)
	}
	return text
}

// returns the text representation of a 'to' message
func (td *textDiagram) toMessageText(message sequencediagram.Message) string {
	var text string
	for i, line := range strings.Split(messageBox(message.Msg), "\n") {
		// add the arrow on the 2nd line
		// length = to_lifeline_index - from_lifeline_index - line_length
		if i == 1 {
			arrowBody := strings.Repeat(arrow_body,
				getPadLength(td.offsets[message.From.Order].getMiddle(), td.offsets[message.To.Order].getMiddle(), line+arrow_start+arrow_forward_end))
			line = arrow_start + line + arrowBody + arrow_forward_end
		} else {
			line = strings.Repeat(" ", utf8.RuneCountInString(arrow_start)) + line
		}
		text += line + "\n"
	}
	return text
}

func (td *textDiagram) fromMessageText(message sequencediagram.Message) string {
	var text string
	msgBox := messageBox(message.Msg)
	// length = from_lifeline_index - to_lifeline_index - line_length
	arrowLength := getPadLength(td.offsets[message.To.Order].getMiddle(), td.offsets[message.From.Order].getMiddle(), arrow_backward_end+arrow_start) - runeIndex(msgBox, '\n')
	for i, line := range strings.Split(msgBox, "\n") {
		// add the arrow on the 2nd line
		if i == 1 {
			arrowBody := strings.Repeat(arrow_body, arrowLength)
			line = arrow_backward_end + arrowBody + line + arrow_start
		} else {
			line = strings.Repeat(" ", arrowLength+utf8.RuneCountInString(arrow_backward_end)) + line
		}
		text += line + "\n"
	}
	return text
}

func (td *textDiagram) drawFullLifeline() {
	var s string
	for i, of := range td.offsets {
		if i == 0 {
			s += strings.Repeat(" ", of.getMiddle()) + life_line
		} else {
			s += strings.Repeat(" ", of.getMiddle()-td.offsets[i-1].getMiddle()-1) + life_line
		}
	}
	td.text += s + "\n"
}

// add lifeline to message text
func (td *textDiagram) fillInLifeline(text string, m sequencediagram.Message) string {
	// toggle lifeline
	defer func() {
		td.lifelineToggle = !td.lifelineToggle
	}()
	if !td.lifelineToggle {
		return text
	}

	// get index start and end nodes
	var startNode, endNode int
	if m.IsSelfLoop() {
		startNode = m.From.Order
		endNode = startNode + 1
		if endNode >= len(td.offsets) {
			endNode = len(td.offsets) - 1
		}
	} else {
		startNode = m.From.Order
		endNode = m.To.Order
		if startNode > endNode {
			startNode, endNode = endNode, startNode
		}
	}
	// use offsets to calculate range of message
	startRange, endRange := td.offsets[startNode].getMiddle(), td.offsets[endNode].getMiddle()
	// for each offset, draw lifeline if it is outside the range of the message
	for _, o := range td.offsets {
		index := o.getMiddle()
		if index <= startRange || index >= endRange {
			if index >= utf8.RuneCountInString(text) {
				text += strings.Repeat(" ", index-utf8.RuneCountInString(text)) + life_line
			} else {
				text = replaceAtRuneIndex(text, index, life_line)
			}
		}
	}
	return text
}