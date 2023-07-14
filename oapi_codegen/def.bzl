load("@io_bazel_rules_go//go:def.bzl", "go_library")

def oapi_codegen_go(name, spec, importpath, visibility, **kwargs):
    codegen_args = {
        "name": name,
        "spec": spec,
    }
    codegen_args.update(kwargs)
    _oapi_codegen_go(**codegen_args)

    go_library(
        name = name + "_generated",
        srcs = [":" + name],
        importpath = importpath,
        deps = [
            "@com_github_deepmap_oapi_codegen//pkg/runtime",
            "@com_github_getkin_kin_openapi//openapi3",
            "@com_github_go_chi_chi_v5//:chi",
        ],
        visibility = visibility,
    )


def _oapi_codegen_go_impl(ctx):
    spec = ctx.file.spec
    output = ctx.actions.declare_file(ctx.label.name + ".gen.go")

    args = ctx.actions.args()
    args.add("-package", ctx.label.name)
    args.add("-o", output)
    args.add("-generate", "types,client,chi-server,strict-server,spec")
    args.add(spec.path)

    ctx.actions.run(
        mnemonic = "APICodegen",
        executable = ctx.executable._codegen,
        arguments = [args],
        inputs = [spec],
        outputs = [output],
    )

    return [
        DefaultInfo(files = depset([output])),
    ]

_oapi_codegen_go = rule(
    implementation = _oapi_codegen_go_impl,
    attrs = {
        "spec": attr.label(
            allow_single_file = True,
            mandatory = True,
        ),
        "_codegen": attr.label(
            default = Label("@com_github_deepmap_oapi_codegen//cmd/oapi-codegen"),
            allow_single_file = True,
            executable = True,
            cfg = "exec",
        ),
    },
)
