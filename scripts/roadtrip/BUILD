load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "roadtrip_lib",
    srcs = ["main.go"],
    importpath = "backend/scripts/roadtrip",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = [
        "//scripts/roadtrip/ai",
        "//scripts/roadtrip/export",
        "//scripts/roadtrip/storage",
        "//scripts/roadtrip/video",
        "@com_github_alecthomas_kong//:kong",
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
