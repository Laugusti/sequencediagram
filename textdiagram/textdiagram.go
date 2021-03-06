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
	title          string
}

// Encode creates an textual representation a sequence diagram using the
// provided sequence diagram
func Encode(sd *sequencediagram.Diagram) io.Reader {
	td := &textDiagram{}
	td.offsets = calcOffsets(sd)
	td.lifelineToggle = true

	nodes := sd.GetOrderedNodes()
	td.addHeaders(nodes, true)
	for _, message := range sd.Messages() {
		td.addMessage(message)
	}
	if td.lifelineToggle {
		td.drawFullLifeline()
	}
	td.addHeaders(nodes, false)
	fixTitle(td)
	return strings.NewReader(td.text)
}

func fixTitle(td *textDiagram) {
	if td.title == "" {
		return
	}
	// split title on "\n"
	title := strings.Split(td.title, "\\n")
	// if there are nodes, center title
	if len(td.offsets) > 0 {
		length := td.offsets[len(td.offsets)-1].end
		for i := range title {
			title[i] = symmetricPadToLength(title[i], ' ', length)
		}
	}
	td.text = strings.Join(title, "\n") + "\n\n" + td.text
}

// addHeaders add the Node slice as text to the ascii diagram
func (td *textDiagram) addHeaders(nodes []*sequencediagram.Node, newline bool) {
	// get max # of lines in Node Name
	height := headerBoxHeight(nodes)
	// headers will have max + 2 lines (top + #lines + bottom)
	headers := make([]string, height+2)
	for i, node := range nodes {
		var pad string
		if i == 0 {
			pad = strings.Repeat(" ", td.offsets[i].begin)
		} else {
			// add padding using pre-calculated node offsets
			pad = strings.Repeat(" ", td.offsets[i].begin-td.offsets[i-1].end-1)
		}
		// add each line of box to header slice with padding
		box := boxString(node.Name, height)
		for j, line := range strings.Split(box, "\n") {
			headers[j] += pad + line
		}
	}
	td.text += strings.Join(headers, "\n")
	if newline {
		td.text += "\n"
	}
}

func (td *textDiagram) paddingForMessage(message sequencediagram.Message) string {
	var pad string
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		pad = strings.Repeat(" ", td.offsets[message.Self.Order].getMiddle()+utf8.RuneCountInString(life_line))
	case sequencediagram.ForwardMessage:
		pad = strings.Repeat(" ", td.offsets[message.From.Order].getMiddle()+utf8.RuneCountInString(life_line))
	case sequencediagram.BackwardMessage:
		pad = strings.Repeat(" ", td.offsets[message.To.Order].getMiddle()+utf8.RuneCountInString(life_line))
	case sequencediagram.Note:
		if message.Side == sequencediagram.Right {
			pad = strings.Repeat(" ", td.offsets[message.Node.Order].getMiddle()+utf8.RuneCountInString(life_line))
		} else {
			if message.Node.Order > 0 {
				box := noteBox(message.Msg)
				length := utf8.RuneCountInString(pad_before_note + box[:strings.Index(box, "\n")])
				pad = strings.Repeat(" ", td.offsets[message.Node.Order].getMiddle()-length)
			}
		}
	}
	return pad
}

// addMessage adds the message the as text to the ascii diagram
func (td *textDiagram) addMessage(message sequencediagram.Message) {
	if t, ok := message.(sequencediagram.Title); ok {
		td.title = t.MessageText()
		return
	}
	text := td.getMessageAsText(message)
	if text == "" {
		return
	}
	// pad with spaces
	pad := td.paddingForMessage(message)

	for _, line := range strings.Split(text, "\n") {
		//for each line, pad and draw life lines
		line = td.fillInLifeline(pad+line, message)
		td.text += line + "\n"
	}
}

// returns the text representation of the message
func (td *textDiagram) getMessageAsText(message sequencediagram.Message) string {
	var text string
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		text = selfLoop(message.Msg, message.AltArrowBody, message.AltArrowEnd)
	case sequencediagram.ForwardMessage:
		text = td.forwardMessageAsText(message)
	case sequencediagram.BackwardMessage:
		text = td.backwardMessageAsText(message)
	case sequencediagram.Note:
		text = noteBox(message.Msg)
	}
	return text
}

// returns the text representation of a 'to' message
func (td *textDiagram) forwardMessageAsText(message sequencediagram.ForwardMessage) string {
	var lines []string
	for i, line := range strings.Split(messageBox(message.Msg), "\n") {
		// add the arrow on the 2nd line
		// length = to_lifeline_index - from_lifeline_index - line_length
		if i == 1 {
			arrowLength := getPadLength(td.offsets[message.From.Order].getMiddle(), td.offsets[message.To.Order].getMiddle(), line+arrow_start+arrow_forward_end)
			line = addArrowToLine(line, arrowLength, message.AltArrowBody, message.AltArrowEnd, false)
		} else {
			line = strings.Repeat(" ", utf8.RuneCountInString(arrow_start)) + line
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func (td *textDiagram) backwardMessageAsText(message sequencediagram.BackwardMessage) string {
	var lines []string
	msgBox := messageBox(message.Msg)
	// length = from_lifeline_index - to_lifeline_index - line_length
	arrowLength := getPadLength(td.offsets[message.To.Order].getMiddle(), td.offsets[message.From.Order].getMiddle(), arrow_backward_end+arrow_start) - runeIndex(msgBox, '\n')
	for i, line := range strings.Split(msgBox, "\n") {
		// add the arrow on the 2nd line
		if i == 1 {
			line = addArrowToLine(line, arrowLength, message.AltArrowBody, message.AltArrowEnd, true)
		} else {
			line = strings.Repeat(" ", arrowLength+utf8.RuneCountInString(arrow_backward_end)) + line
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// add an arrow to the line
func addArrowToLine(line string, arrowLength int, altArrowBody, altArrowEnd, backwards bool) string {
	arrowStart := arrow_start
	arrowBody := strings.Repeat(arrow_body, arrowLength)
	if altArrowBody {
		arrowStart = alt_arrow_start
		arrowBody = strings.Repeat(alt_arrow_body, arrowLength)
	}
	if backwards {
		arrowEnd := arrow_backward_end
		if altArrowEnd {
			arrowEnd = alt_arrow_backward_end
		}
		line = arrowEnd + arrowBody + line + arrowStart
	} else {
		arrowEnd := arrow_forward_end
		if altArrowEnd {
			arrowEnd = alt_arrow_forward_end
		}
		line = arrowStart + line + arrowBody + arrowEnd
	}
	return line
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
func (td *textDiagram) fillInLifeline(text string, message sequencediagram.Message) string {
	// toggle lifeline
	defer func() {
		td.lifelineToggle = !td.lifelineToggle
	}()
	if !td.lifelineToggle {
		return text
	}

	startRange, endRange := td.getStartEndIndex(message)
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

// getStartEndIndex returns the start and end index of the Message
func (td *textDiagram) getStartEndIndex(message sequencediagram.Message) (int, int) {
	var startNode, endNode int
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		startNode = message.Self.Order
		endNode = startNode + 1
		if endNode >= len(td.offsets) {
			endNode = len(td.offsets) - 1
		}
	case sequencediagram.ForwardMessage:
		startNode, endNode = message.From.Order, message.To.Order
	case sequencediagram.BackwardMessage:
		startNode, endNode = message.To.Order, message.From.Order
	}

	// use offsets to calculate range of message
	return td.offsets[startNode].getMiddle(), td.offsets[endNode].getMiddle()
}
