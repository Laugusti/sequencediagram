package sequencediagram

import "fmt"

type Message interface {
	MessageText() string
	String() string
}

type noMessage struct{}

func (nm noMessage) MessageText() string { return "" }

type simpleMessage struct {
	Msg string
}

func (sm simpleMessage) MessageText() string {
	return sm.Msg
}

type uniDirectionalMessage struct {
	AltArrowBody bool
	AltArrowEnd  bool
}

func (udm uniDirectionalMessage) arrow() string {
	arrowBody := "-"
	arrowEnd := ">"
	if udm.AltArrowBody {
		arrowBody = "--"
	}
	if udm.AltArrowEnd {
		arrowEnd = ">>"
	}
	return arrowBody + arrowEnd
}

type Title struct {
	simpleMessage
}

func (t Title) String() string {
	return "title " + t.Msg
}

type Participant struct {
	Self *Node
	noMessage
}

func (p Participant) String() string {
	return "participant " + p.Self.Name
}

type SelfMessage struct {
	Self *Node
	simpleMessage
	uniDirectionalMessage
}

func (sm SelfMessage) String() string {
	return fmt.Sprintf("%s%s%s:%s", sm.Self.Name, sm.arrow(), sm.Self.Name, sm.Msg)
}

type ForwardMessage struct {
	From *Node
	To   *Node
	simpleMessage
	uniDirectionalMessage
}

func (fm ForwardMessage) String() string {
	return fmt.Sprintf("%s%s%s:%s", fm.From.Name, fm.arrow(), fm.To.Name, fm.Msg)
}

type BackwardMessage struct {
	From *Node
	To   *Node
	simpleMessage
	uniDirectionalMessage
}

func (bm BackwardMessage) String() string {
	return fmt.Sprintf("%s%s%s:%s", bm.From.Name, bm.arrow(), bm.To.Name, bm.Msg)
}
