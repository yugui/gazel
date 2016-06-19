package generator_test

import (
	"go/build"
	"os"
	"path/filepath"
	"testing"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/yugui/gazel/generator"
)

func testData() string {
	srcdir := os.Getenv("TEST_SRCDIR")
	if srcdir == "" {
		return "testdata"
	}
	return filepath.Join(srcdir, os.Getenv("TEST_WORKSPACE"), "generator", "testdata")
}

func canonicalize(t *testing.T, filename, content string) string {
	f, err := bzl.Parse(filename, []byte(content))
	if err != nil {
		t.Fatalf("bzl.Parse(%q, %q) failed with %v; want success", filename, content, err)
	}
	return string(bzl.Format(f))
}

func format(rules []*bzl.Rule) string {
	var f bzl.File
	for _, r := range rules {
		f.Stmt = append(f.Stmt, r.Call)
	}
	return string(bzl.Format(&f))
}

func packageFromDir(t *testing.T, dir string) *build.Package {
	pkg, err := build.ImportDir(dir, build.ImportComment)
	if err != nil {
		t.Fatalf("build.ImportDir(%q, build.ImportComment) failed with %v; want success", dir, err)
	}
	return pkg
}

func TestGeneratorWithLib(t *testing.T) {
	g := generator.New("example.com/repo")
	pkg := packageFromDir(t, filepath.Join(testData(), "lib"))
	rules, err := g.Generate(".", pkg)
	if err != nil {
		t.Errorf(`g.Generate(".", %#v) failed with %v; want success`, pkg, err)
	}

	want := canonicalize(t, "BUILD", `
		go_library(
			name = "go_default_library",
			srcs = ["doc.go", "lib.go"],
		)

		go_test(
			name = "go_default_test",
			srcs = ["lib_test.go"],
			library = ":go_default_library",
		)
	`)
	if got := format(rules); got != want {
		t.Errorf(`g.Generate(".", %#v) = %s; want %s`, pkg, got, want)
	}
}

func TestGeneratorWithSubdirLib(t *testing.T) {
	g := generator.New("example.com/repo")
	pkg := packageFromDir(t, filepath.Join(testData(), "lib"))
	rules, err := g.Generate("lib", pkg)
	if err != nil {
		t.Errorf(`g.Generate("lib", %#v) failed with %v; want success`, pkg, err)
	}

	want := canonicalize(t, "lib/BUILD", `
		go_library(
			name = "lib",
			srcs = ["doc.go", "lib.go"],
		)

		go_test(
			name = "lib_test",
			srcs = ["lib_test.go"],
			library = ":lib",
		)
	`)
	if got := format(rules); got != want {
		t.Errorf(`g.Generate(".", %#v) = %s; want %s`, pkg, got, want)
	}
}

func TestGeneratorWithBin(t *testing.T) {
	g := generator.New("example.com/repo")
	pkg := packageFromDir(t, filepath.Join(testData(), "bin"))
	rules, err := g.Generate("bin", pkg)
	if err != nil {
		t.Errorf(`g.Generate("bin", %#v) failed with %v; want success`, pkg, err)
	}

	want := canonicalize(t, "bin/BUILD", `
		go_binary(
			name = "bin",
			srcs = ["main.go"],
			deps = ["//lib:go_default_library"],
		)
	`)
	if got := format(rules); got != want {
		t.Errorf(`g.Generate(".", %#v) = %s; want %s`, pkg, got, want)
	}
}
