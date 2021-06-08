package span

import (
	"roc/internal/x"
)

// attention:it's maybe zipkin instead,don't user this for your production Environment.

// Span
// |---TraceId:1     ----->RPC----->       |---TraceId:1
//       |---ParentSpanId:0                            |---ParentSpanId:222
//           |---SpanId:222                                 |---SpanId:223
type Span struct {
	SpanId       uint32
	ParentSpanId int32
	traceId      string
}

func (s *Span) Carrier() {
	s.ParentSpanId += 1
	s.SpanId += 1
}

func (s *Span) Finish() {
	// todo buffer flush to cloud or something
	return
}

func (s *Span) String() string {
	return "span"
}

func (s *Span) TraceId() string {
	return s.traceId
}

func NewSpan() *Span {
	return &Span{
		traceId:      x.NewUUID(),
		ParentSpanId: -1,
		SpanId:       1,
	}
}
