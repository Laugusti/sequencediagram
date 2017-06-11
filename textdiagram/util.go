package textdiagram

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Laugusti/sequencediagram"
)

func symmetricPadToLength(s string, r rune, n int) string {
	if len(s) >= n {
		return s
	}
	padLeft := strings.Repeat(string(r), (n-len(s))/2)
	padRight := strings.Repeat(string(r), (n-len(s)+1)/2)
	return padLeft + s + padRight
}

// replaces the ith rune in the s with r
func replaceAtRuneIndex(s string, i int, new string) string {
	var runeSize int
	var runeCount int
	for idx := 0; idx < len(s); idx += runeSize {
		_, runeSize = utf8.DecodeRuneInString(s[idx:])
		if runeCount == i {
			s = fmt.Sprintf("%s%s%s", s[:idx], new, s[idx+runeSize:])
			break
		}
		runeCount++
	}
	return s
}

// finds the rune index of r (# of runes before r) in s, returns -1 if not found
func runeIndex(s string, r rune) int {
	var runeCount int
	for i := 0; i < len(s); {
		r2, size := utf8.DecodeRuneInString(s[i:])
		if r == r2 {
			return runeCount
		}
		i += size
		runeCount++
	}
	return -1
}

func getPadLength(startIndex, endIndex int, otherCharacters string) int {
	otherLength := utf8.RuneCountInString(otherCharacters)
	return endIndex - startIndex - 1 - otherLength
}

func headerBoxHeight(nodes []*sequencediagram.Node) int {
	var max int
	for _, node := range nodes {
		lines := len(strings.Split(node.Name, "\\n"))
		if lines > max {
			max = lines
		}
	}
	return max
}
