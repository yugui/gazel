// Command gazel is a BUILD file generator for Go projects.
// See "gazel --help" for more details.
package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/yugui/gazel/generator"
)

var (
	goPrefix = flag.String("go_prefix", "", "go_prefix of the target workspace")
)

func generate(bctx build.Context, goPrefix, root string) error {
	var rules []bzl.Expr
	g := generator.New(goPrefix)
	err := generator.Walk(bctx, root, func(dir string, pkg *build.Package) error {
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

func run(goPrefix string, dirs []string) error {
	bctx := build.Default
	// Ignore $GOPATH environment variable
	bctx.GOPATH = ""
	for _, d := range dirs {
		if err := generate(bctx, goPrefix, d); err != nil {
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

	if err := run(*goPrefix, flag.Args()); err != nil {
		log.Fatal(err)
	}
}
