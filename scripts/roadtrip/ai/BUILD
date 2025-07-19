load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ai",
    srcs = ["ai.go"],
    importpath = "backend/scripts/roadtrip/ai",
    visibility = ["//scripts/roadtrip:__subpackages__"],
    deps = [
        "@com_github_google_generative_ai_go//genai:genai",
        "@org_golang_google_api//option:option",
    ],
)

go_test(
    name = "ai_test",
    srcs = ["ai_test.go"],
    data = glob(["testdata/**"]) + glob(["../prompts/**"]),
    embed = [":ai"],
)
