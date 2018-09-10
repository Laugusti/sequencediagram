package textdiagram

import "testing"

func TestBoxString(t *testing.T) {
	tests := []struct {
		text   string
		height int
		want   string
	}{
		{``, 1, "┌──┐\n│  │\n└──┘"},
		{`abc`, 1, "┌─────┐\n│ abc │\n└─────┘"},
		{`abc`, 2, "┌─────┐\n│ abc │\n│     │\n└─────┘"},
		{`abc`, 3, "┌─────┐\n│ abc │\n│     │\n│     │\n└─────┘"},
		{`abc\ndefghi\nj`, 1, "┌────────┐\n│  abc   │\n│ defghi │\n│   j    │\n└────────┘"},
	}

	for _, test := range tests {
		got := boxString(test.text, test.height)
		if got != test.want {
			t.Errorf("TestBoxString => got wrong box: input (text: %q height: %d), got: %q, want: %q", test.text, test.height, got, test.want)
		}
	}
}

func TestSelfLoop(t *testing.T) {
	tests := []struct {
		text    string
		altBody bool
		altEnd  bool
		want    string
	}{
		{``, false, false, "────┐\n    │ \n◀───┘"},
		{``, true, false, "----┐\n    ¦ \n◀---┘"},
		{``, false, true, "────┐\n    │ \n<───┘"},
		{``, true, true, "----┐\n    ¦ \n<---┘"},
		{`abc`, false, false, "────┐\n    │abc \n◀───┘"},
		{`a\nbcd\nef`, false, false, "────┐\n    │a \n    │bcd \n    │ef \n◀───┘"},
	}

	for _, test := range tests {
		got := selfLoop(test.text, test.altBody, test.altEnd)
		if got != test.want {
			t.Errorf("TestSelfLoop => got wrong loop: input (text: %q altbody: %t altend: %t), got: %q, want: %q", test.text, test.altBody, test.altEnd, got, test.want)
		}
	}
}

func TestMessageBox(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{``, "┌──┐\n┤  ├\n└──┘"},
		{`abc`, "┌─────┐\n┤ abc ├\n└─────┘"},
		{`abc\ndefghi\nj`, "┌────────┐\n┤  abc   ├\n│ defghi │\n│   j    │\n└────────┘"},
	}

	for _, test := range tests {
		got := messageBox(test.text)
		if got != test.want {
			t.Errorf("TestMessageBox => got wrong box: input (text: %q), got: %q, want: %q", test.text, got, test.want)
		}
	}
}

func TestNoteBox(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{``, " ┌──╗\n │  │\n └──┘"},
		{`abc`, " ┌─────╗\n │ abc │\n └─────┘"},
		{`abc\ndefghi\nj`, " ┌────────╗\n │  abc   │\n │ defghi │\n │   j    │\n └────────┘"},
	}

	for _, test := range tests {
		got := noteBox(test.text)
		if got != test.want {
			t.Errorf("TestNoteBox => got wrong box: input (text: %q), got: %q, want: %q", test.text, got, test.want)
		}
	}
}
