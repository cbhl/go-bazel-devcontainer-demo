load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ai",
    srcs = ["ai.go"],
    importpath = "backend/scripts/roadtrip/ai",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_test(
    name = "ai_test",
    srcs = ["ai_test.go"],
    data = glob(["testdata/**"]),
    embed = [":ai"],
)
