# Response

[![codecov](https://codecov.io/gh/cantor-systems/response/branch/master/graph/badge.svg?token=bWx6yzysUm)](https://codecov.io/gh/cantor-systems/response)

Use it:

```
response.With(w, r, http.StatusOK, data)
```

Package respond provides low-touch API responses for Go data services.

  * Idiomatic way of responding to data APIs using `response.With`
  * Use `response.With` to respond with default options, or make a `response.Options` for [advanced features](https://godoc.org/go.cantor.systems/response/respond#Options)
  * Encoder abstraction lets you easily speak different formats
  * `Before` and `After` function fields allow you to envelope and mutate data, set common HTTP headers, log activity etc.
  * Protected against multiple responses

