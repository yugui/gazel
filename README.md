gazel
======

gazel is a BUILD file generator for Go projects.

## Status
Prototype

## Features
* Generates `go_library` for library packages
* Generates `go_binary` for command packages
* Collects `srcs`
* Collects `deps` within the same repository

## TODO
* Generate `go_test`
* Collect `deps` from external dependency
* Collect `deps` from vendor directory
* Build tags, release tags
* Support generating a `BUILD` file for each Go package
* cgo
* SWIG
  * once [`rules_go`](https://github.com/bazelbuild/rules_go) supports SWIG.
* Respect manually configured existing rules
