package sequencediagram

import "testing"

func testParseFromText(t *testing.T) {
	tests := []struct {
		text        string
		shouldParse bool
	}{
		{"test", false},
		{"alice->bob:msg", true},
		{"alice->bob:msg\na->b:c", true},
		{"alice->bob:msg\n", false},
	}
	for _, test := range tests {
		sd, err := ParseFromText(test.text)
		if test.shouldParse && err != nil {
			t.Errorf("expected no parse errors but got err = %v", err)
		} else if !test.shouldParse && err == nil {
			t.Errorf("expected parse errors")
		} else if sd != nil && sd.String() != test.text {
			t.Errorf("expected %q got %q", test.text, sd)
		}
	}
}
