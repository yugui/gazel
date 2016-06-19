package main

import (
	"io/ioutil"
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
)

func fixFile(fname string, rules []bzl.Expr) (err error) {
	buildfile, err := reconcile(fname, rules)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "BUILD")
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			if merr := os.Rename(f.Name(), fname); merr != nil {
				err = merr
			}
		}
	}()
	defer func() {
		if cerr := f.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()

	_, err = f.Write(bzl.Format(buildfile))
	return err
}
