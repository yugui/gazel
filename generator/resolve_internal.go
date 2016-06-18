package generator

import (
	"fmt"
	"strings"
)

// internalResolver resolves go_library labels within the same repository as
// the one of goPrefix.
type internalResolver struct {
	goPrefix string
}

// resolve takes a Go importpath within the same respository as r.goPrefix
// and resolves it into a label in Bazel.
func (r internalResolver) resolve(importpath string) (label, error) {
	if importpath == r.goPrefix {
		return label{name: "go_default_library"}, nil
	}

	if prefix := r.goPrefix + "/"; strings.HasPrefix(importpath, prefix) {
		return label{
			pkg:  strings.TrimPrefix(importpath, prefix),
			name: "go_default_library",
		}, nil
	}

	return label{}, fmt.Errorf("importpath %q does not start with goPrefix %q", importpath, r.goPrefix)
}
