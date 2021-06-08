package simple

import (
	"roc/internal/x"
)

// |---TraceId:1     ----->RPC----->       |---TraceId:1
type Simple struct {
	traceId string
}

func (s *Simple) Carrier() {
	return
}

func NewSimple(traceId ...string) *Simple {
	s := &Simple{}
	if len(traceId) > 0 && traceId[0] != "" {
		s.traceId = traceId[0]
		return s
	}
	s.traceId = x.NewUUID()
	return s
}

func (s *Simple) String() string {
	return "simple"
}

func (s *Simple) Finish() {
	return
}

func (s *Simple) TraceId() string {
	return s.traceId
}
