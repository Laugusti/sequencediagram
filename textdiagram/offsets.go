package textdiagram

import (
	"strings"
	"unicode/utf8"

	"github.com/Laugusti/sequencediagram"
)

type offset struct {
	begin int
	end   int
}

func (o offset) getMiddle() int {
	return o.begin + (o.end-o.begin+1)/2
}

// calOffsets create an offset for each node in the sequence diagram, the order
// of the node is the index in the offset slice
func calcOffsets(sd *sequencediagram.Diagram) []offset {
	nodes := sd.GetOrderedNodes()
	offsets := make([]offset, len(nodes))
	// calc minimum offsets between nodes
	for i, node := range nodes {
		// begin index is 0 for 1st node or the index of the last node + 1
		var begin int
		if i > 0 {
			begin = offsets[i-1].end + 1
		}
		// end index is begin + number of runes in box - 1
		boxSize := utf8.RuneCountInString(box_vertical+node.Name+box_vertical) + 2*box_inside_pad
		end := begin + boxSize - 1
		offsets[i] = offset{begin, end}
	}

	// adjust offsets based on message
	for _, message := range sd.Messages() {
		// calculate begining node index to start shifting. do nothing if shift is past last node
		shiftStart := calcShiftStartIndex(message)
		if shiftStart >= len(nodes) {
			continue
		}

		// calculate required shift, do nothing if shift is not required
		shift := calcShift(message, offsets)
		if shift < 1 {
			continue
		}

		// shift offsets
		for i := shiftStart; i < len(offsets); i++ {
			offsets[i].begin += shift
			offsets[i].end += shift
		}
	}
	return offsets
}

// calcShiftStartIndex calculates from which node to start shifting offsets based on the message
func calcShiftStartIndex(message sequencediagram.Message) int {
	var shiftStart int
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		shiftStart = message.Self.Order + 1
	case sequencediagram.ForwardMessage:
		shiftStart = message.To.Order
	case sequencediagram.BackwardMessage:
		shiftStart = message.From.Order
	}
	return shiftStart
}

// calcShift calculates required shift to the offset based on the message
func calcShift(message sequencediagram.Message, offsets []offset) int {
	// get length of longest string in message
	var length int
	for _, m := range splitMessage(message) {
		if utf8.RuneCountInString(m) > length {
			length = utf8.RuneCountInString(m)
		}
	}

	// calculate shift based on message
	var shift int
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		length += utf8.RuneCountInString(loop_middle + pad_between_loop_and_message + loop_message_end_pad)
		offset1 := offsets[message.Self.Order].getMiddle()
		offset2 := offsets[message.Self.Order+1].getMiddle()
		diff := offset2 - offset1 - 1
		if length > diff {
			shift = length - diff
		}
	case sequencediagram.ForwardMessage:
		length += utf8.RuneCountInString(arrow_start+box_arrow_left+box_arrow_right+arrow_body+arrow_forward_end) + 2*box_inside_pad
		offset1 := offsets[message.From.Order].getMiddle()
		offset2 := offsets[message.To.Order].getMiddle()
		diff := offset2 - offset1 - 1
		if length > diff {
			shift = length - diff
		}
	case sequencediagram.BackwardMessage:
		length += utf8.RuneCountInString(arrow_backward_end+arrow_body+box_arrow_left+box_arrow_right+arrow_start) + 2*box_inside_pad
		offset1 := offsets[message.To.Order].getMiddle()
		offset2 := offsets[message.From.Order].getMiddle()
		diff := offset2 - offset1 - 1
		if length > diff {
			shift = length - diff
		}
	}
	return shift
}

func getMessageText(message sequencediagram.Message) string {
	switch message := message.(type) {
	case sequencediagram.SelfMessage:
		return message.Msg
	case sequencediagram.ForwardMessage:
		return message.Msg
	case sequencediagram.BackwardMessage:
		return message.Msg
	}
	return ""
}

func splitMessage(message sequencediagram.Message) []string {
	return strings.Split(getMessageText(message), "\\n")
}
