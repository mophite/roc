package trace

type Trace interface {
	Carrier()
	Finish()
	TraceId() string
	String() string
}