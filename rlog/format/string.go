package format

import (
	"bytes"

	"roc/internal/x/bytesbuffpool"
	"roc/rlog/common"
)

var _ Formatter = &stringFormat{}

type stringFormat struct {
	layout string
}

func (s *stringFormat) Layout() string {
	if s.layout == "" {
		return defaultLayout
	}
	return s.layout
}

func (s *stringFormat) Format(detail *common.Detail) *bytes.Buffer {
	b := bytesbuffpool.Get()
	b.WriteString("[" + detail.Level + "] ")
	b.WriteString(detail.Timestamp + " ")
	if detail.Line != "" {
		b.WriteString(detail.Line + " ")
	}
	b.WriteString(detail.Content)
	return b
}

func (s *stringFormat) SetLayout(layout string) {
	s.layout = layout
}

func (s *stringFormat) String() string {
	return "str"
}
