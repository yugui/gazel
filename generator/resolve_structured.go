package generator

import (
	"fmt"
	"strings"
)

// structuredResolver resolves go_library labels within the same repository as
// the one of goPrefix.
type structuredResolver struct {
	goPrefix string
}

// resolve takes a Go importpath within the same respository as r.goPrefix
// and resolves it into a label in Bazel.
func (r structuredResolver) resolve(importpath, dir string) (label, error) {
	if importpath == r.goPrefix {
		return label{name: "go_default_library"}, nil
	}

	if prefix := r.goPrefix + "/"; strings.HasPrefix(importpath, prefix) {
		pkg := strings.TrimPrefix(importpath, prefix)
		if pkg == dir {
			return label{name: "go_default_library", relative: true}, nil
		}
		return label{pkg: pkg, name: "go_default_library"}, nil
	}

	return label{}, fmt.Errorf("importpath %q does not start with goPrefix %q", importpath, r.goPrefix)
}
