package generator

import (
	"reflect"
	"testing"
)

func TestInternalResolver(t *testing.T) {
	r := internalResolver{goPrefix: "example.com/repo"}

	for _, spec := range []struct {
		importpath string
		want       label
	}{
		{
			importpath: "example.com/repo",
			want:       label{name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/sub",
			want:       label{pkg: "sub", name: "go_default_library"},
		},
	} {
		l, err := r.resolve(spec.importpath)
		if err != nil {
			t.Errorf("r.resolve(%q) failed with %v; want success", spec.importpath, err)
			continue
		}
		if got, want := l, spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("r.resolve(%q) = %s; want %s", spec.importpath, got, want)
		}
	}
}

func TestInternalResolverError(t *testing.T) {
	r := internalResolver{goPrefix: "example.com/repo"}

	for _, importpath := range []string{
		"example.com/another",
		"example.com/another/sub",
		"example.com/repo_suffix",
	} {
		l, err := r.resolve(importpath)
		if err == nil {
			t.Errorf("r.resolve(%q) = %s; want error", importpath, l)
		}
	}
}
