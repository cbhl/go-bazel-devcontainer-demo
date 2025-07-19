load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "export",
    srcs = ["export.go"],
    importpath = "backend/scripts/roadtrip/export",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_test(
    name = "export_test",
    srcs = ["export_test.go"],
    data = glob(["testdata/**"]),
    embed = [":export"],
)
