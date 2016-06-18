package generator

import (
	"go/build"
	"os"
	"path/filepath"
)

// A WalkFunc is a callback called by Walk for each package.
// The first parameter "dir" is a relative path to the package directory from the given root.
// "dir" is "." for the root itself.
// The secnd parameter "pkg" is package metadata.
type WalkFunc func(dir string, pkg *build.Package) error

// Walk walks through Go packages under the given dir.
// It calls back "f" for each package.
func Walk(bctx build.Context, root string, f WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		pkg, err := bctx.ImportDir(path, build.ImportComment)
		if _, ok := err.(*build.NoGoError); ok {
			return nil
		}
		if err != nil {
			return err
		}
		return f(rel, pkg)
	})
}
