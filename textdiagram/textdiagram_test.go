package textdiagram

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/Laugusti/sequencediagram"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{readFile(t, "testdata/test1_sd.txt"), readFile(t, "testdata/test1_td.txt")},
		{readFile(t, "testdata/test2_sd.txt"), readFile(t, "testdata/test2_td.txt")},
	}
	for _, test := range tests {
		got := getAsTextDiagram(t, test.text)
		if got != test.want {
			t.Errorf("TestEncode => input: %q\n, got:\n%q\n\twant:\n%q", test.text, got, test.want)
		}
	}
}

func readFile(t *testing.T, filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	return string(bytes.TrimSpace(b))
}

func getAsTextDiagram(t *testing.T, s string) string {
	sd, err := sequencediagram.ParseFromText(s)
	if err != nil {
		t.Fatalf("error parsing sequence diagram: %v", err)
	}
	r := Encode(sd)
	var b bytes.Buffer
	io.Copy(&b, r)
	return b.String()
}
