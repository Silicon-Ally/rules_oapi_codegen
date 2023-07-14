# Example Usage

This minimal example demonstrates the funcionality of the `oapi_codegen_go` build rule.

To run this you'll need to have Bazel installed, and it's also recommended to use Bazelisk to choose a suitable version of Bazel to run. We use Gazelle to manage Go dependencies in Bazel, but that isn't strictly necessary.
## Running the example

Run the server:

```bash
$ bazel run //cmd/server -- --port=8080
```

Make some API calls:

```bash
$ curl localhost:8080/pets

[]
```

```bash
$ curl \
  -X POST \
  --data '{"name": "Scruffles", "tag": "good dog"}' \
  -H 'Content-Type: application/json' \
  localhost:8080/pets

{"id":1,"name":"Scruffles","tag":"good dog"}
```

```bash
$ curl localhost:8080/pets

[{"id":1,"name":"Scruffles","tag":"good dog"}]
```
