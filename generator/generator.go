package generator

import (
	"fmt"
	"go/build"
	"path"
	"path/filepath"
	"strings"

	bzl "github.com/bazelbuild/buildifier/core"
)

// Generator generates Bazel build rules for a Go package.
type Generator interface {
	// Generate generates build rules for a Go package.
	// "dir" is a relative path from the repository root to the directory of
	// the Go package.
	// "pkg" is a description about the package.
	Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error)
}

// A Mode describes how Generator organizes rules for different Go packages.
type Mode int

const (
	// FlatMode means that Generator puts all rules for different Go packages
	// in the current repository into a single top level package of Bazel.
	FlatMode = Mode(iota) + 1
	// StructuredMode means that Generator generates a Bazel package for each
	// Go package.
	StructuredMode
)

// New returns an implementation of Generator.
// "goPrefix" is the go_prefix corresponding to the repository root.
// "mode" specifies how to organize rules for different Go packages.
func New(goPrefix string, mode Mode) Generator {
	var r0 labelResolver
	switch mode {
	case FlatMode:
		r0 = flatResolver{goPrefix: goPrefix}
	case StructuredMode:
		r0 = structuredResolver{goPrefix: goPrefix}
	default:
		panic(fmt.Sprintf("unrecognized mode %d", mode))
	}

	var e externalResolver
	r := resolverFunc(func(importpath, dir string) (label, error) {
		if importpath != goPrefix && !strings.HasPrefix(importpath, goPrefix+"/") && !strings.HasPrefix(importpath, "./") {
			return e.resolve(importpath, dir)
		}
		return r0.resolve(importpath, dir)
	})
	return &generator{
		goPrefix: goPrefix,
		r:        r,
	}
}

type generator struct {
	goPrefix string
	r        labelResolver
}

func (g *generator) Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error) {
	r, err := g.generate(filepath.Base(pkg.Dir), dir, pkg.GoFiles, pkg.Imports, pkg.IsCommand())
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

	if len(pkg.XTestGoFiles) > 0 {
		t, err := g.generateXTest(dir, pkg.XTestGoFiles, pkg.XTestImports)
		if err != nil {
			return nil, err
		}
		rules = append(rules, t)
	}

	return rules, nil
}

func (g *generator) generate(basename, rel string, srcs, imports []string, isCommand bool) (*bzl.Rule, error) {
	l, err := g.r.resolve(path.Join(g.goPrefix, rel), rel)
	if err != nil {
		return nil, err
	}
	name := l.name

	kind := "go_library"
	if isCommand {
		kind = "go_binary"
		name = basename
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: srcs},
	}

	deps, err := g.dependencies(imports, rel)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}

	return newRule(kind, nil, attrs)
}

func (g *generator) generateTest(dir string, srcs, imports []string, library string) (*bzl.Rule, error) {
	l, err := g.r.resolve(path.Join(g.goPrefix, dir), dir)
	if err != nil {
		return nil, err
	}
	name := l.name + "_test"
	if l.name == "go_default_library" {
		name = "go_default_test"
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: srcs},
		{key: "library", value: ":" + library},
	}

	deps, err := g.dependencies(imports, dir)
	if err != nil {
		return nil, err
	}
	if len(deps) > 0 {
		attrs = append(attrs, keyvalue{key: "deps", value: deps})
	}
	return newRule("go_test", nil, attrs)
}

func (g *generator) generateXTest(dir string, srcs, imports []string) (*bzl.Rule, error) {
	l, err := g.r.resolve(path.Join(g.goPrefix, dir), dir)
	if err != nil {
		return nil, err
	}
	name := l.name + "_xtest"
	if l.name == "go_default_library" {
		name = "go_default_xtest"
	}

	attrs := []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: srcs},
	}

	deps, err := g.dependencies(imports, dir)
	if err != nil {
		return nil, err
	}
	attrs = append(attrs, keyvalue{key: "deps", value: deps})
	return newRule("go_test", nil, attrs)
}

func (g *generator) dependencies(imports []string, dir string) ([]string, error) {
	var deps []string
	for _, p := range imports {
		if isStandard(p) {
			continue
		}
		l, err := g.r.resolve(p, dir)
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
