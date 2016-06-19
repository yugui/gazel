package generator

import (
	"reflect"
	"testing"
)

func TestStructuredResolver(t *testing.T) {
	r := structuredResolver{goPrefix: "example.com/repo"}
	for _, spec := range []struct {
		importpath string
		curPkg     string
		want       label
	}{
		{
			importpath: "example.com/repo",
			curPkg:     "",
			want:       label{name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/lib",
			curPkg:     "",
			want:       label{pkg: "lib", name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/another",
			curPkg:     "",
			want:       label{pkg: "another", name: "go_default_library"},
		},

		{
			importpath: "example.com/repo",
			curPkg:     "lib",
			want:       label{name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/lib",
			curPkg:     "lib",
			want:       label{name: "go_default_library", relative: true},
		},
		{
			importpath: "example.com/repo/lib/sub",
			curPkg:     "lib",
			want:       label{pkg: "lib/sub", name: "go_default_library"},
		},
		{
			importpath: "example.com/repo/another",
			curPkg:     "lib",
			want:       label{pkg: "another", name: "go_default_library"},
		},
	} {

		l, err := r.resolve(spec.importpath, spec.curPkg)
		if err != nil {
			t.Errorf("r.resolve(%q) failed with %v; want success", spec.importpath, err)
			continue
		}
		if got, want := l, spec.want; !reflect.DeepEqual(got, want) {
			t.Errorf("r.resolve(%q) = %s; want %s", spec.importpath, got, want)
		}
	}
}

func TestStructuredResolverError(t *testing.T) {
	r := structuredResolver{goPrefix: "example.com/repo"}

	for _, importpath := range []string{
		"example.com/another",
		"example.com/another/sub",
		"example.com/repo_suffix",
	} {
		l, err := r.resolve(importpath, "")
		if err == nil {
			t.Errorf("r.resolve(%q) = %s; want error", importpath, l)
		}
	}
}
