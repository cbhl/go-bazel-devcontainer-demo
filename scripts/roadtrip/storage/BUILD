load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "storage",
    srcs = ["storage.go"],
    importpath = "backend/scripts/roadtrip/storage",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = [
        "@com_github_minio_minio_go_v7//:minio-go",
        "@org_golang_google_cloud_storage//:storage",
    ],
)

go_test(
    name = "storage_test",
    srcs = ["storage_test.go"],
    embed = [":storage"],
    data = glob(["testdata/**"]),
)