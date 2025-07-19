load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "env",
    srcs = [
        "config.go",
        "env.go",
    ],
    importpath = "backend/scripts/roadtrip/env",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_test(
    name = "env_test",
    srcs = ["env_test.go"],
    embed = [":env"],
)
