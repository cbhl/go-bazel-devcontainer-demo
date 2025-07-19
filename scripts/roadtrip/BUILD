load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "roadtrip_lib",
    srcs = ["main.go"],
    deps = [
        "//scripts/roadtrip/env:env",
        "@com_github_alecthomas_kong//:kong",
    ],
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_binary(
    name = "roadtrip",
    srcs = ["main.go"],
    deps = [
        "//scripts/roadtrip/env:env",
        "@com_github_alecthomas_kong//:kong",
    ],
    visibility = ["//visibility:public"],
)

go_test(
    name = "roadtrip_test",
    srcs = ["main_test.go"],
    embed = [":roadtrip_lib"],
    deps = [
        "//scripts/roadtrip/env:env",
    ],
)

