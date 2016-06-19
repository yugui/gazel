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

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/yugui/gazel/generator"
)

var (
	goPrefix = flag.String("go_prefix", "", "go_prefix of the target workspace")
	flat     = flag.Bool("flat", false, "creates a large single BUILD file in the top of repository instead of creating a BUILD file for each Go package")
)

func generate(bctx build.Context, mode generator.Mode, root string) error {
	g := generator.New(*goPrefix, mode)

	drive := func(bctx build.Context, root string, f generator.WalkFunc) error {
		pkg, err := bctx.ImportDir(root, build.ImportComment)
		if err != nil {
			return err
		}
		return f("", pkg)
	}
	if filepath.Base(root) == "..." {
		drive = generator.Walk
		root = filepath.Dir(root)
	}

	var rules []bzl.Expr
	err := drive(bctx, root, func(dir string, pkg *build.Package) error {
		rs, err := g.Generate(dir, pkg)
		if err != nil {
			return err
		}
		for _, r := range rs {
			rules = append(rules, r.Call)
		}
		return nil
	})
	if err != nil {
		return err
	}
	buf := bzl.Format(&bzl.File{
		Stmt: rules,
	})
	_, err = os.Stdout.Write(buf)
	return err
}

func run(dirs []string) error {
	bctx := build.Default
	// Ignore $GOPATH environment variable
	bctx.GOPATH = ""

	m := generator.StructuredMode
	if *flat {
		m = generator.FlatMode
	}

	for _, d := range dirs {
		if err := generate(bctx, m, d); err != nil {
			return err
		}
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, `Gazel is a BUILD file generator fo Go projects.

Currently its primary usage is to generate BUILD files for external dependencies
in a go_vendor repository rule.
You can still use Gazel for other purposes, but its interface can change without notice.
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
	if len(flag.Args()) > 1 && *flat {
		log.Fatal("can have only one argument when -flat=true")
	}

	if err := run(flag.Args()); err != nil {
		log.Fatal(err)
	}
}
