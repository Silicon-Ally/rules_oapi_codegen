# rules_oapi_codegen

[![CI Workflow](https://github.com/Silicon-Ally/rules_oapi_codegen/actions/workflows/build.yml/badge.svg)](https://github.com/Silicon-Ally/rules_oapi_codegen/actions?query=branch%3Amain)

`rules_oapi_codegen` provides [Bazel](https://bazel.build) rules for generating Go code from OpenAPI 3.0 definitions. It uses github.com/deepmap/oapi-codegen to provide this functionality.

## Usage

Look at [the Releases page](https://github.com/Silicon-Ally/rules_oapi_codegen/releases) for instructions on how to update your Bazel `WORKSPACE`.

```bazel
# In a BUILD.bazel file

load("@rules_oapi_codegen//oapi_codegen:def.bzl", "oapi_codegen_go")

oapi_codegen_go(
    name = "api",
    importpath = "<your import path>",
    spec = "<your OpenAPI YAML file>",
    visibility = ["//visibility:public"],
)
```

## Example

The `example` directory provides [the standard Petstore example](https://github.com/OAI/OpenAPI-Specification/blob/main/examples/vv3.0/petstore.yaml), see [the example/ README](/example/README.md) for more details.

## Contributing

Contribution guidelines can be found [on our website](https://siliconally.org/oss/contributor-guidelines).