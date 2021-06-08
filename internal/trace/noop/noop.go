package noop

type Noop struct{}

func (n *Noop) Carrier() {
	return
}

func NewSimple() *Noop {
	return new(Noop)
}

func (n *Noop) String() string {
	return "noop"
}

func (n *Noop) Finish() {
	return
}

func (n *Noop) TraceId() string {
	return ""
}
