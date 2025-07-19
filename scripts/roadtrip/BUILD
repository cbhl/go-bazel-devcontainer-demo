load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "roadtrip_lib",
    srcs = ["main.go"],
    importpath = "backend/scripts/roadtrip",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = [
        "@com_github_alecthomas_kong//:kong",
        "//scripts/roadtrip/video",
    ],
)

go_binary(
    name = "roadtrip",
    embed = [":roadtrip_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "roadtrip_test",
    srcs = ["main_test.go"],
    data = glob(["testdata/**"]),
    embed = [":roadtrip_lib"],
)
