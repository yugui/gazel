package generator

import (
	"fmt"
	"go/build"
	"path/filepath"
	"reflect"

	bzl "github.com/bazelbuild/buildifier/core"
)

// Generator generates Bazel build rules for a Go package.
type Generator interface {
	// Generate generates build rules for a Go package.
	// "dir" is a relative path from the current Bazel package directory to the Go package.
	// "pkg" is a description about the package.
	Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error)
}

// New returns an implementation of Generator.
func New() Generator {
	return new(generator)
}

type generator struct{}

func (g *generator) Generate(dir string, pkg *build.Package) ([]*bzl.Rule, error) {
	kind := "go_library"
	name := filepath.ToSlash(dir)
	if pkg.IsCommand() {
		kind = "go_binary"
		name = filepath.Base(dir)
	}
	if dir == "." {
		name = "go_default_library"
	}

	var rules []*bzl.Rule
	r, err := newRule(kind, nil, []keyvalue{
		{key: "name", value: name},
		{key: "srcs", value: pkg.GoFiles},
	})
	if err != nil {
		return nil, err
	}
	rules = append(rules, r)
	return rules, nil
}

type keyvalue struct {
	key   string
	value interface{}
}

func newRule(kind string, args []interface{}, kwargs []keyvalue) (*bzl.Rule, error) {
	var list []bzl.Expr
	for i, arg := range args {
		expr, err := newValue(arg)
		if err != nil {
			return nil, fmt.Errorf("wrong arg %v at args[%d]: %v", arg, i, err)
		}
		list = append(list, expr)
	}
	for _, arg := range kwargs {
		expr, err := newValue(arg.value)
		if err != nil {
			return nil, fmt.Errorf("wrong value %v at kwargs[%q]: %v", arg.value, arg.key, err)
		}
		list = append(list, &bzl.BinaryExpr{
			X:  &bzl.LiteralExpr{Token: arg.key},
			Op: "=",
			Y:  expr,
		})
	}

	return &bzl.Rule{
		Call: &bzl.CallExpr{
			X:    &bzl.LiteralExpr{Token: kind},
			List: list,
		},
	}, nil
}

// newValue converts a Go value into the corresponding expression in Bazel BUILD file.
func newValue(val interface{}) (bzl.Expr, error) {
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &bzl.LiteralExpr{Token: fmt.Sprintf("%d", val)}, nil
	case reflect.Float32, reflect.Float64:
		return &bzl.LiteralExpr{Token: fmt.Sprintf("%f", val)}, nil
	case reflect.String:
		return &bzl.StringExpr{Value: val.(string)}, nil
	case reflect.Slice, reflect.Array:
		var list []bzl.Expr
		for i := 0; i < rv.Len(); i++ {
			elem, err := newValue(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			list = append(list, elem)
		}
		return &bzl.ListExpr{List: list}, nil
	default:
		return nil, fmt.Errorf("not implemented %T", val)
	}
}
