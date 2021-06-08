package namespace

import (
	"strings"
)

// SplicingConfigPrefix is splicing config prefix
// Schema default is DefaultConfigSchema
func SplicingConfigPrefix(schema Schema, prefix string) string {
	var b strings.Builder
	b.WriteString(schema)

	if prefix != "" {
		b.WriteString("/")
		b.WriteString(prefix)
	}

	return b.String()
}
