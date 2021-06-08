package format

import (
	"bytes"

	"roc/rlog/common"
)

const defaultLayout = "2006.01.02.15:04:05.000"

var DefaultFormat Formatter = &stringFormat{}

type Formatter interface {
	Layout() string
	Format(detail *common.Detail) *bytes.Buffer
	SetLayout(layout string)
	String() string
}
