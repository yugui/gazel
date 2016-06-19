load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "gazel",
    srcs = [
        "diff.go",
        "fix.go",
        "main.go",
        "print.go",
        "reconcile.go",
    ],
    deps = [
        "@io_bazel_buildifier//core:go_default_library",
        "@io_bazel_buildifier//differ:go_default_library",
        "//generator:go_default_library",
    ],
)
