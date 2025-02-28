# Named Primitive Type References Example

This example demonstrates how to reuse primitive-based types in the OpenAPI spec. 

By default, Huma doesn't create references (`$ref`) for custom types based on primitives (like `type CustomHeader string`). Instead, it defines them inline at each usage point. This keeps the OpenAPI specification simpler.

However, in some cases you might want these types to be defined once and referenced via `$ref` elsewhere:
- When you have custom primitive-based types used in multiple places
- When you need consistent schema documentation for those types
- When you want smaller API specifications by avoiding redundant schema definitions

## How it works

The example defines:
1. Two custom primitive-based types (`CustomHeader` and `CustomQueryParam`)
2. Two request structs that use these types
3. An API that enables the `ReuseNamedPrimitiveTypes` option

```go
// Enable the reuse of primitive-based types
config := huma.DefaultConfig("My API", "1.0.0")
config.ReuseNamedPrimitiveTypes = true
```

## Running the example

```bash
go run .
```

Then open:
- API documentation: http://localhost:8888/docs
- OpenAPI specification: http://localhost:8888/openapi.json

## What to look for

In the OpenAPI specification, observe that the `CustomHeader` and `CustomQueryParam` types are:
1. Defined once in the `components/schemas` section
2. Referenced via `$ref` in the request body schemas

Without this option enabled, the types would be defined inline in each request schema instead of using references. 