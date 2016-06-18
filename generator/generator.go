package generator

import (
	"go/build"
	"path/filepath"
	"strings"

	bzl "github.com/bazelbuild/buildifier/core"
)

// Generator generates Bazel build rules for a Go package.
type Generator interface {
	// Generate generates build rules for a Go package.
	// "dir" is a relative path from the current Bazel package directory to the Go package.
	// "pkg" is a description about the package.
	Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error)
}

// New returns an implementation of Generator.
func New(goPrefix string) Generator {
	return &generator{
		r: internalResolver{
			goPrefix: goPrefix,
		},
	}
}

type generator struct {
	r labelResolver
}

func (g *generator) Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error) {
	kind := "go_library"
	name := filepath.ToSlash(dir)
	if pkg.IsCommand() {
		kind = "go_binary"
		name = filepath.Base(dir)
	}
	if dir == "." {
		name = "go_default_library"
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: pkg.GoFiles},
	}

	deps, err := g.dependencies(pkg)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}

	var rules []*bzl.Rule
	r, err := newRule(kind, nil, attrs)
	if err != nil {
		return nil, err
	}
	rules = append(rules, r)
	return rules, nil
}

func (g *generator) dependencies(pkg *build.Package) ([]string, error) {
	var deps []string
	for _, p := range pkg.Imports {
		if g.isStandard(p) {
			continue
		}
		l, err := g.r.resolve(p)
		if err != nil {
			return nil, err
		}
		deps = append(deps, l.String())
	}
	return deps, nil
}

// isStandard determines if importpath points a Go standard package.
func (g *generator) isStandard(importpath string) bool {
	seg := strings.SplitN(importpath, "/", 2)[0]
	return !strings.Contains(seg, ".")
}
