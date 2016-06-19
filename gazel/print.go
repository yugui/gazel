package main

import (
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
)

func printFile(fname string, rules []bzl.Expr) (err error) {
	buildfile, err := reconcile(fname, rules)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(bzl.Format(buildfile))
	return err
}
