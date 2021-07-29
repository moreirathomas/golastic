# `pkg/golastic`

The package `golastic` defines a custom API to interface `elastic/go-elasticsearch` package.

It exposes 3 main APIs: `Indices`, `Document` and `Search`. Each of these have methods to interact with Elasticsearch. These methods return custom response types which themselves have method to conveniently parse and use the native `esapi.Response` response type.

> `golastic` options and structs are defined based on our needs. Thus requests are not fully parameterizable nor the response types are exhaustive regarding the fields defined and supported.

## Make a request

You must retrieve the corresponding `golastic` API and provide a context for the request.

Then simply chain call the method for the request you are making.

```go
// Indices API
res, _ := golastic.Indices(ctx).CreateIfNotExists("my-index")

// Document API
res, _ := golastic.Document(ctx).Index(doc)

// Search API
res, _ := golastic.Search(ctx).MultiMatchQuery("foo", fields, pagination, sort)
```

## Use the response

Each `golastic` API methods return their own response type.

Most of the time, you access the data by calling `Unwrap()` on the response. A destination struct may be provided to unmarshal the response as configured.

For more details on each response type and their convenience methods, refer to the source code.

```go
result, _ := res.Unwrap(MyStruct{}) // result is interface{}

t, ok := result.(MyStruct) // t is MyStruct
```
