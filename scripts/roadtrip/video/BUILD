load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "video",
    srcs = ["video.go"],
    importpath = "backend/scripts/roadtrip/video",
    visibility = ["//scripts/roadtrip:__subpackages__"],
)

go_test(
    name = "video_test",
    srcs = ["video_test.go"],
    data = ["//scripts/roadtrip:testdata"],
    embed = [":video"],
)
