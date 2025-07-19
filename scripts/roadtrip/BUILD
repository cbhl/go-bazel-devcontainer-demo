load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "roadtrip",
    srcs = ["main.go"],
    deps = [
        "//scripts/roadtrip/env",
        "@com_github_alecthomas_kong//:kong",
    ],
    visibility = ["//visibility:public"],
)

go_library(
    name = "env",
    srcs = [
        "env/env.go",
        "env/config.go",
    ],
    importpath = "backend/scripts/roadtrip/env",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

# Build test to verify the binary builds
build_test(
    name = "roadtrip_build_test",
    targets = [":roadtrip"],
)