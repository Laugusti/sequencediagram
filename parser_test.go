package sequencediagram

import "reflect"
import "testing"

func TestParseFromText(t *testing.T) {
	tests := []struct {
		text        string
		shouldParse bool
	}{
		{"test", false},
		{"title", false},
		{"title title", true},
		{"title title\ntitle title", true},
		{"participant alice", true},
		{"title title\nparticipant alice", true},
		{"alice->alice:msg", true},
		{"alice->alice:multi\\nline\\nmsg", true},
		{"alice->bob:msg", true},
		{"alice->bob:msg\na->b:c", true},
		{"a->b:req\\nb->a:resp", true},
		{"a-->b:msg", true},
		{"a->>b:msg", true},
		{"a-->>b:msg", true},
		{"alice->bob:msg\n", false},
		{"note right of alice:msg", true},
		{"note left of alice:msg", true},
		{"note above alice:msg", false},
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

func TestParseFromTextMessageType(t *testing.T) {
	tests := []struct {
		text         string
		messageTypes []reflect.Type
	}{
		{"title title", messageTypes(Title{})},
		{"participant a", messageTypes(Participant{})},
		{"title title\nparticipant a", messageTypes(Title{}, Participant{})},
		{"a->a:msg", messageTypes(SelfMessage{})},
		{"a->b:msg", messageTypes(ForwardMessage{})},
		{"a->b:msg\nb-->>a:resp", messageTypes(ForwardMessage{}, BackwardMessage{})},
		{"note right of a:msg", messageTypes(Note{})},
		{"note left of a:msg", messageTypes(Note{})},
	}

	for _, test := range tests {
		sd, err := ParseFromText(test.text)
		if err != nil {
			t.Errorf("TestParseFromTextMessageType => got parse error: %v", err)
			continue
		}
		if len(sd.Messages()) != len(test.messageTypes) {
			t.Errorf("TestParseFromTextMessageType => expected %d messages, got %d", len(test.messageTypes), len(sd.messages))
			continue
		}
		for i, gotMsg := range sd.Messages() {
			if reflect.TypeOf(gotMsg) != test.messageTypes[i] {
				t.Errorf("TestParseFromTextMessageType => expected %v message type, got %T", test.messageTypes[i], gotMsg)
			}
		}
	}
}

func messageTypes(messages ...Message) (types []reflect.Type) {
	for _, msg := range messages {
		types = append(types, reflect.TypeOf(msg))
	}
	return
}
