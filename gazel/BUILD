load("@io_bazel_rules_go//go:def.bzl", "go_binary")

go_binary(
    name = "gazel",
    srcs = ["main.go"],
    deps = [
        "//generator:go_default_library",
        "@io_bazel_buildifier//core:go_default_library",
    ],
)
