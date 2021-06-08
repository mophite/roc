package format

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"

	"roc/internal/x/bytesbuffpool"
	"roc/rlog/common"
)

var _ Formatter = &jsonFormat{}

type jsonFormat struct {
	layout string
}

func (j *jsonFormat) Layout() string {
	if j.layout == "" {
		return defaultLayout
	}
	return j.layout
}

func (j *jsonFormat) Format(detail *common.Detail) *bytes.Buffer {
	b := bytesbuffpool.Get()
	b.Write(mustMarshal(detail))
	return b
}

func (j *jsonFormat) SetLayout(layout string) {
	j.layout = layout
}

func (j *jsonFormat) String() string {
	return "json"
}

var fastest = jsoniter.ConfigFastest

func mustMarshal(v interface{}) []byte {
	b, _ := fastest.Marshal(v)
	return b
}
