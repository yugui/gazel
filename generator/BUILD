load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "walk.go",
        "generator.go",
    ],
    deps = [
        "@io_bazel_buildifier//core:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "walk_test.go",
        "generator_test.go",
    ],
    deps = [
        ":go_default_library",
        "@io_bazel_buildifier//core:go_default_library",
    ], 
    data = glob(["testdata/**/*"]),
)
