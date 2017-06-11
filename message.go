package sequencediagram

type Message struct {
	From *Node
	To   *Node
	Msg  string
}

// TODO: change Message to interface and add implementation for
// SelfMessage, ToMessage, FromMessage
type SelfMessage struct {
	Self *Node
	Msg  string
}

func (m Message) IsSelfLoop() bool {
	return m.From == m.To
}
