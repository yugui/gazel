// Command gazel is a BUILD file generator for Go projects.
// See "gazel --help" for more details.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path/filepath"
	"strings"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/yugui/gazel/generator"
)

var (
	goPrefix = flag.String("go_prefix", "", "go_prefix of the target workspace")
	baseDir  = flag.String("base_dir", "", "path to a directory which corresponds to go_prefix")
	flat     = flag.Bool("flat", false, "creates a large single BUILD file in the top of repository instead of creating a BUILD file for each Go package")
	mode     = flag.String("mode", "print", "print, fix or diff")
)

type gen struct {
	base string
	bctx build.Context
	g    generator.Generator
	emit func(fname string, rules []bzl.Expr) error
}

func newGen() (*gen, error) {
	base, err := filepath.Abs(*baseDir)
	if err != nil {
		return nil, err
	}

	bctx := build.Default
	// Ignore $GOPATH environment variable
	bctx.GOPATH = ""

	m := generator.StructuredMode
	if *flat {
		m = generator.FlatMode
	}

	g := gen{
		base: filepath.Clean(base),
		bctx: bctx,
		g:    generator.New(*goPrefix, m),
	}
	switch *mode {
	case "print":
		g.emit = printFile
	case "fix":
		g.emit = fixFile
	case "diff":
		g.emit = diffFile
	}
	return &g, nil
}

func (g gen) generate(root string) error {
	drive := func(bctx build.Context, root string, f generator.WalkFunc) error {
		pkg, err := bctx.ImportDir(root, build.ImportComment)
		if err != nil {
			return err
		}
		return f(pkg)
	}
	if filepath.Base(root) == "..." {
		drive = generator.Walk
		root = filepath.Dir(root)
	}

	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	root = filepath.Clean(root)
	if root != g.base && !strings.HasPrefix(root, fmt.Sprintf("%s%c", g.base, filepath.Separator)) {
		return fmt.Errorf("dir %s is not under the base dir %s", root, g.base)
	}

	var rules []bzl.Expr
	if root == g.base {
		rules = append(rules, &bzl.CallExpr{
			X: &bzl.LiteralExpr{Token: "go_prefix"},
			List: []bzl.Expr{
				&bzl.StringExpr{Value: *goPrefix},
			},
		})
	}
	err = drive(g.bctx, root, func(pkg *build.Package) error {
		rel, err := filepath.Rel(g.base, pkg.Dir)
		if err != nil {
			return err
		}
		if rel == "." {
			rel = ""
		}

		rs, err := g.g.Generate(filepath.ToSlash(rel), pkg)
		if err != nil {
			return err
		}
		for _, r := range rs {
			rules = append(rules, r.Call)
		}

		if !*flat {
			if err := g.emit(filepath.Join(pkg.Dir, "BUILD"), rules); err != nil {
				return err
			}
			rules = nil
		}
		return nil
	})
	if err != nil {
		return err
	}

	if *flat {
		return g.emit(filepath.Join(root, "BUILD"), rules)
	}
	return nil
}

func run(dirs []string) error {
	g, err := newGen()
	if err != nil {
		return err
	}

	for _, d := range dirs {
		if err := g.generate(d); err != nil {
			return err
		}
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: gazel [flags...] [package-dirs...]

Gazel is a BUILD file generator for Go projects.

Currently its primary usage is to generate BUILD files for external dependencies
in a go_vendor repository rule.
You can still use Gazel for other purposes, but its interface can change without
notice.

It takes a list of paths to Go package directories.
It recursively traverses its subpackages if the directory path ends with "/...".
All the directories must be under the directory specified in -base_dir.

There are several modes of gazel.
In print mode, gazel prints reconciled BUILD files to stdout.
In fix mode, gazel creates BUILD files or updates existing ones.
In diff mode, gazel shows diff.

FLAGS:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *goPrefix == "" {
		// TODO(yugui) Extract go_prefix from the top level BUILD file if exists
		log.Fatal("-go_prefix is required")
	}
	if *baseDir == "" {
		if flag.NArg() != 1 {
			log.Fatal("-base_dir is required")
		}
		*baseDir = flag.Arg(0)
		if dir, base := filepath.Split(*baseDir); base == "..." {
			*baseDir = dir
		}
	}
	if len(flag.Args()) > 1 && *flat {
		log.Fatal("can have only one argument when -flat=true")
	}

	if err := run(flag.Args()); err != nil {
		log.Fatal(err)
	}
}
