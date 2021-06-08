package namespace

import (
	"strings"
)

// SplicingPrefix is splicing service address
// Schema default is DefaultSchema. eg. goroc
func SplicingPrefix(schema Schema, scope Scope) string {
	var b strings.Builder
	b.WriteString(schema)

	if scope != "" {
		b.WriteString("/")
		b.WriteString(scope)
	}

	return b.String()
}

// SplicingScope is the service name/version. eg.srv.hello/version
// name or version don't allowed none
func SplicingScope(name, version string) Scope {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString("/")
	b.WriteString(version)
	return b.String()
}
