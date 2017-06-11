package textdiagram

import (
	"strings"
	"unicode/utf8"
)

const (
	box_vertical   = "│"
	box_horizontal = "─"

	box_top_left     = "┌"
	box_top_right    = "┐"
	box_bottom_left  = "└"
	box_bottom_right = "┘"
	box_arrow_left   = "┤"
	box_arrow_right  = "├"

	arrow_start        = "──"
	arrow_body         = "─"
	arrow_forward_end  = "▶"
	arrow_backward_end = "◀"
	arrow_vertical     = "│"

	alt_arrow_start        = "--"
	alt_arrow_body         = "-"
	alt_arrow_forward_end  = ">"
	alt_arrow_backward_end = "<"
	alt_arrow_vertical     = "¦"

	life_line     = "│"
	alt_life_line = "‖"
)

const (
	box_inside_pad = 1

	loop_body_length             = 3
	pad_between_loop_and_message = ""
	loop_message_end_pad         = " "
)

// boxString wraps s in a text box, padding to padToHeight if necessary
func boxString(s string, padToHeight int) string {
	lines := strings.Split(s, "\\n")

	// get max line length
	var maxLength int
	for _, line := range lines {
		if utf8.RuneCountInString(line) > maxLength {
			maxLength = utf8.RuneCountInString(line)
		}
	}
	maxLength += 2 * box_inside_pad
	// pad lines and wrap in box verticals
	var content string
	for _, line := range lines {
		content += box_vertical + symmetricPadToLength(line, ' ', maxLength) + box_vertical + "\n"
	}
	// pad height if necessary
	for i := len(lines); i < padToHeight; i++ {
		content += box_vertical + strings.Repeat(box_horizontal, maxLength) + box_vertical + "\n"
	}

	// create box
	middle := strings.Repeat(box_horizontal, maxLength)
	box := box_top_left + middle + box_top_right + "\n"
	box += content
	box += box_bottom_left + middle + box_bottom_right
	return box
}

// selfLoop text diagram of a arrow that loops back to self with message s
func selfLoop(s string, altArrowBody, altArrowEnd bool) string {
	loopTop := strings.Repeat(arrow_body, loop_body_length+1) + box_top_right
	loopMiddle := strings.Repeat(" ", loop_body_length+1) + arrow_vertical
	loopBottom := arrow_backward_end + strings.Repeat(arrow_body, loop_body_length) + box_bottom_right
	if altArrowEnd {
		loopBottom = strings.Replace(loopBottom, arrow_backward_end, alt_arrow_backward_end, 1)
	}
	if altArrowBody {
		loopTop = strings.Replace(loopTop, arrow_body, alt_arrow_body, -1)
		loopMiddle = strings.Replace(loopMiddle, arrow_vertical, alt_arrow_vertical, 1)
		loopBottom = strings.Replace(loopBottom, arrow_body, alt_arrow_body, -1)

	}
	loop := loopTop + "\n"
	for _, line := range strings.Split(s, "\\n") {
		loop += loopMiddle + pad_between_loop_and_message + line + loop_message_end_pad + "\n"
	}
	loop += loopBottom
	return loop
}

// messageBox is similar boxString except for the walls of the 2nd line
func messageBox(s string) string {
	box := strings.Split(boxString(s, 0), "\n")
	for i, line := range box {
		if i == 1 {
			line = replaceAtRuneIndex(line, utf8.RuneCountInString(line)-1, box_arrow_right)
			box[i] = replaceAtRuneIndex(line, 0, box_arrow_left)
		}
	}
	return strings.Join(box, "\n")
}
