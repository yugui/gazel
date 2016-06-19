package main

import (
	"io/ioutil"
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
)

func reconcile(fname string, rules []bzl.Expr) (*bzl.File, error) {
	buf, err := ioutil.ReadFile(fname)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	orig := new(bzl.File)
	if len(buf) > 0 {
		orig, err = bzl.Parse(fname, buf)
		if err != nil {
			return nil, err
		}
	}

	newfile := *orig
	// TODO(yugui) Respect existing data, visibility and other attributes;
	// comments on rules; and their positions.
	newfile.DelRules("go_library", "")
	newfile.DelRules("go_binary", "")
	newfile.DelRules("go_test", "")
	newfile.Stmt = append(newfile.Stmt, rules...)
	return &newfile, nil
}
