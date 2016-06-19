load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "walk.go",
        "construct.go",
        "generator.go",
        "resolve.go",
        "resolve_flat.go",
        "resolve_structured.go",
    ],
    deps = [
        "@io_bazel_buildifier//core:go_default_library",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "generator_external_test",
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

go_test(
    name = "generator_test",
    srcs = [
        "resolve_flat_test.go",
        "resolve_structured_test.go",
    ],
    library = ":go_default_library",
    data = glob(["testdata/**/*"]),
)
