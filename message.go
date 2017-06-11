package sequencediagram

import "fmt"

type Message interface {
	MessageText() string
	String() string
}

type simpleMessage struct {
	Msg string
}

func (sm simpleMessage) MessageText() string {
	return sm.Msg
}

type SelfMessage struct {
	Self *Node
	simpleMessage
}

func (sm SelfMessage) String() string {
	return fmt.Sprintf("%s->%s:%s", sm.Self.Name, sm.Self.Name, sm.Msg)
}

type ForwardMessage struct {
	From *Node
	To   *Node
	simpleMessage
}

func (fm ForwardMessage) String() string {
	return fmt.Sprintf("%s->%s:%s", fm.From.Name, fm.To.Name, fm.Msg)
}

type BackwardMessage struct {
	From *Node
	To   *Node
	simpleMessage
}

func (bm BackwardMessage) String() string {
	return fmt.Sprintf("%s->%s:%s", bm.From.Name, bm.To.Name, bm.Msg)
}
