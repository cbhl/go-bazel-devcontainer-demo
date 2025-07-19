load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "storage",
    srcs = ["storage.go"],
    importpath = "backend/scripts/roadtrip/storage",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = [
        "@com_google_cloud_go_storage//:storage",
    ],
)

go_test(
    name = "storage_test",
    srcs = ["storage_test.go"],
    embed = [":storage"],
)