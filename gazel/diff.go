package main

import (
	"io/ioutil"
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/bazelbuild/buildifier/differ"
)

func diffFile(fname string, rules []bzl.Expr) (err error) {
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
			if merr := os.Remove(f.Name()); merr != nil {
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
	if _, err := f.Write(bzl.Format(buildfile)); err != nil {
		return err
	}

	diff := differ.Find()
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		diff.Show(os.DevNull, f.Name())
		return nil
	}
	diff.Show(fname, f.Name())
	return nil
}
