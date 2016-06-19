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
	r, err := g.generate(dir, pkg.GoFiles, pkg.Imports, pkg.IsCommand())
	if err != nil {
		return nil, err
	}
	rules := []*bzl.Rule{r}

	if len(pkg.TestGoFiles) > 0 {
		t, err := g.generateTest(dir, pkg.TestGoFiles, pkg.TestImports, r.AttrString("name"))
		if err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}

	return rules, nil
}

func (g *generator) generate(dir string, srcs, imports []string, isCommand bool) (*bzl.Rule, error) {
	kind := "go_library"
	name := filepath.ToSlash(dir)
	if isCommand {
		kind = "go_binary"
		name = filepath.Base(dir)
	}
	if dir == "." {
		name = "go_default_library"
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: srcs},
	}

	deps, err := g.dependencies(imports)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}

	return newRule(kind, nil, attrs)
}

func (g *generator) generateTest(dir string, srcs, imports []string, library string) (*bzl.Rule, error) {
	name := filepath.ToSlash(dir) + "_test"
	if dir == "." {
		name = "go_default_test"
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: srcs},
		{key: "library", value: ":" + library},
	}

	deps, err := g.dependencies(imports)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}
	return newRule("go_test", nil, attrs)
}

func (g *generator) dependencies(imports []string) ([]string, error) {
	var deps []string
	for _, p := range imports {
		if isStandard(p) {
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
func isStandard(importpath string) bool {
	seg := strings.SplitN(importpath, "/", 2)[0]
	return !strings.Contains(seg, ".")
}
