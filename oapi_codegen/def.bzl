load("@io_bazel_rules_go//go:def.bzl", "go_library")

def oapi_codegen_go(name, spec, importpath, visibility, skip_prune=False, **kwargs):
    """Generates Go bindings for an OpenAPI 3.0 spec.

    This rule runs [oapi-codegen](https://github.com/deepmap/oapi-codegen) to
    produce a generated Go library providing, among other things, a strictly
    typed API stub interface to implement your API against.

    Args:
      name: A unique name for this rule.
      spec: The OpenAPI 3.0 YAML specification for your API, must be self-contained.
      importpath: The importpath of the directory the rule is defined in, like
        'github.com/<org>/<repo>/path/to/dir/api'. This is the import path of
        the generated Go library
      visibility: The visibility of the generated go_library target.
      skip_prune: Whether or not to generate code for unreferenced schema components.
    """

    codegen_args = {
        "name": name,
        "spec": spec,
        "skip_prune": skip_prune,
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

    targets = ["types","client","chi-server","strict-server","spec"]
    if ctx.attr.skip_prune:
        targets.append("skip-prune")

    args = ctx.actions.args()
    args.add("-package", ctx.label.name)
    args.add("-o", output)
    args.add("-generate", ",".join(targets))
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
        "skip_prune": attr.bool(
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
