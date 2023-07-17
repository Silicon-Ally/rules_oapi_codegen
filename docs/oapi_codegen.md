<!-- Generated with Stardoc: http://skydoc.bazel.build -->



<a id="oapi_codegen_go"></a>

## oapi_codegen_go

<pre>
oapi_codegen_go(<a href="#oapi_codegen_go-name">name</a>, <a href="#oapi_codegen_go-spec">spec</a>, <a href="#oapi_codegen_go-importpath">importpath</a>, <a href="#oapi_codegen_go-visibility">visibility</a>, <a href="#oapi_codegen_go-kwargs">kwargs</a>)
</pre>

Generates Go bindings for an OpenAPI 3.0 spec.

This rule runs [oapi-codegen](https://github.com/deepmap/oapi-codegen) to
produce a generated Go library providing, among other things, a strictly
typed API stub interface to implement your API against.


**PARAMETERS**


| Name  | Description | Default Value |
| :------------- | :------------- | :------------- |
| <a id="oapi_codegen_go-name"></a>name |  A unique name for this rule.   |  none |
| <a id="oapi_codegen_go-spec"></a>spec |  The OpenAPI 3.0 YAML specification for your API, must be self-contained.   |  none |
| <a id="oapi_codegen_go-importpath"></a>importpath |  The importpath of the directory the rule is defined in, like 'github.com/&lt;org&gt;/&lt;repo&gt;/path/to/dir/api'. This is the import path of the generated Go library   |  none |
| <a id="oapi_codegen_go-visibility"></a>visibility |  The visibility of the generated go_library target.   |  none |
| <a id="oapi_codegen_go-kwargs"></a>kwargs |  <p align="center"> - </p>   |  none |


