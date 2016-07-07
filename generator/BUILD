load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "construct.go",
        "generator.go",
        "resolve.go",
        "resolve_external.go",
        "resolve_flat.go",
        "resolve_structured.go",
        "walk.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@io_bazel_buildifier//core:go_default_library",
    ],
)

go_test(
    name = "generator_external_test",
    srcs = [
        "generator_test.go",
        "walk_test.go",
    ],
    data = glob(["testdata/**/*"]),
    deps = [
        "@io_bazel_buildifier//core:go_default_library",
        ":go_default_library",
    ],
)

go_test(
    name = "generator_test",
    srcs = [
        "resolve_flat_test.go",
        "resolve_structured_test.go",
    ],
    data = glob(["testdata/**/*"]),
    library = ":go_default_library",
)
