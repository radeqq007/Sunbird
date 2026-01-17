# http

`http` is a module that provides functions for handling HTTP requests.

```ts
import "http"
```

## create_server

`create_server` is a function that creates a new HTTP server.

```ts
const server = http.create_server()
```
`server` is an object that provides methods for handling HTTP requests.

### server route methods

`server` has the following route methods:
- `get`
- `post`
- `put`
- `delete`
- `patch`
- `options`
- `connect`
- `trace`

They all work the same way, only the HTTP method is different.

```ts
server.get("/", func(w, r) {
    w.send("Hello World!")
})

server.post("/", func(w, r) {
    w.send("Hello World!")
})

server.put("/", func(w, r) {
    w.send("Hello World!")
})

server.delete("/", func(w, r) {
    w.send("Hello World!")
})

server.patch("/", func(w, r) {
    w.send("Hello World!")
})

server.options("/", func(w, r) {
    w.send("Hello World!")
})

server.head("/", func(w, r) {
    w.send("Hello World!")
})
```

All of them take two arguments:
- `route`: the route to register the handler for
- `callback`: the callback function to handle the request

`callback` takes two arguments:
- `w`: the response writer
- `r`: the request

To see more about the writer and the request, see the [writer](#writer) and [request](#request) sections.


You can register a path parameter by using `{name}` in the route.

```ts
server.get("/user/{id}", func(w, r) {
    w.send("Hello World!")
})
```

The path parameters will be available in the request object.

```ts
const id = r.path_param("id")
```
### listen

`listen` is a server method that starts the server and takes an integer as an argument.

```ts
server.listen(8080)
```

This will start an HTTP server on port 8080.

## writer

`writer` is an object that provides methods for writing to the response.
It gets passed as the first argument to the route handler.


It has the following methods:
- `send`: sends a string to the response
- `json`: sends a JSON object to the response
- `cookie`: object that contains functions for handling cookies
- `header`: object that provides methods for setting, adding, deleting, and getting headers
- `status`: sets the status code of the response

### writer.send

`writer.send` is a function that sends a string to the response.

```ts
server.get("/", func(w, r) {
    w.send("Hello World!")
})
```

### writer.json

`writer.json` is a function that sends a JSON object to the response.

```ts
server.get("/", func(w, r) {
    w.json({ "message": "Hello World!" })
})
```

### writer.cookie

`writer.cookie` is an object that contains functions for handling cookies.

#### writer.cookie.set()

`writer.cookie.set()` is a function that sets a cookie.

It takes two arguments:
- `name`: the name of the cookie
- `value`: the value of the cookie

and a third optional argument:
- `options`: an object that can contain: `max_age`, `domain`, `path`, `secure`, `http_only`, `same_site`

```ts
server.get("/", func(w, r) {
    w.set_cookie("name", "value")
})
```

```ts
server.get("/", func(w, r) {
    w.cookie.set("name", "value", {
        "max_age": 3600,         // Cookie expiration time in seconds
        "domain": "example.com", // Cookie domain
        "path": "/",             // Cookie path
        "secure": true,          // Only send over HTTPS
        "http_only": true,       // Not accessible via JavaScript
        "same_site": "strict"    // CSRF protection ("strict", "lax", or "none")
    })
})
```

#### writer.cookie.delete()

`writer.cookie.delete` is a function that deletes a cookie.

```ts
server.get("/", func(w, r) {
    w.cookie.delete("name")
})
```

### writer.status

`writer.status` is a function that sets the status code of the response.

```ts
server.get("/", func(w, r) {
    w.status(http.status.ok)
})
```

To see all the status codes, see the [status](#status) section.

### writer.header

`writer.header` is an object that provides methods for setting, adding, deleting, and getting headers.

#### writer.header.set

`writer.header.set` is a function that sets a header.

```ts
server.get("/", func(w, r) {
    w.header.set("Content-Type", "application/json")
})
```

#### `writer.header.add`

`writer.header.add` is a function that adds a header.

```ts
server.get("/", func(w, r) {
    w.header.add("Content-Type", "application/json")
})
```

#### `writer.header.del`

`writer.header.del` is a function that deletes a header.

```ts
server.get("/", func(w, r) {
    w.header.del("Content-Type")
})
```

#### `writer.header.get`

`writer.header.get` is a function that gets a header value.

```ts
server.get("/", func(w, r) {
    const contentType = w.header.get("Content-Type")
})
```

## request

`request` is an object that provides methods for getting information about the request.
It gets passed as the second argument to the route handler.

It contains the following methods:
- `path_param`: gets a path parameter
- `query_param`: gets a query parameter
- `body`: gets the request body
- `json`: gets the request body as a JSON object
- `method`: gets the request method
- `url`: gets the request URL
- `header`: gets the request header value
- `headers`: gets the request headers
- `cookie`: gets the request cookie value
- `cookies`: gets the request cookies

### request.path_param

`request.path_param` is a function that gets a path parameter.

```ts
server.get("/user/{id}", func(w, r) {
    const id = r.path_param("id")
})
```

### request.query_param

`request.query_param` is a function that gets a query parameter.

```ts
server.get("/user", func(w, r) {
    const id = r.query_param("id")
})
```

When the query parameter is not found, it returns `null`.

### request.body

`request.body` is a function that gets the request body.

```ts
server.post("/user", func(w, r) {
    const body = r.body()
})
```

### request.json

`request.json` is a function that gets the request body as a JSON object.

```ts
server.post("/user", func(w, r) {
    const body = r.json()
})
```

### request.method

`request.method` is a function that gets the request method.

```ts
server.get("/user", func(w, r) {
    const method = r.method()
})
```

### request.url

`request.url` is a function that gets the request URL.

```ts
server.get("/user", func(w, r) {
    const url = r.url()
})
```

### request.header

`request.header` is a function that gets the request header value.

```ts
server.get("/user", func(w, r) {
    const contentType = r.header("Content-Type")
})
```

When the header is not found, it returns `null`.

### request.headers

`request.headers` is a function that gets the request headers.

```ts
server.get("/user", func(w, r) {
    const headers = r.headers()
})
```

### request.cookie

`request.cookie` is a function that gets the request cookie value.

```ts
server.get("/user", func(w, r) {
    const cookie = r.cookie("name")
})
```

If the cookie is not found, it returns `null`.

### request.cookies

`request.cookies` is a function that gets the request cookies.

```ts
server.get("/user", func(w, r) {
    const cookies = r.cookies()
})
```

## status

`status` is an object that provides constants for HTTP status codes.

```ts
http.status.continue // 100
http.status.switching_protocols // 101
http.status.processing // 102
http.status.early_hints // 103
http.status.ok // 200
http.status.created // 201
http.status.accepted // 202
http.status.non_authoritative_info // 203
http.status.no_content // 204
http.status.reset_content // 205
http.status.partial_content // 206
http.status.multi_status // 207
http.status.already_reported // 208
http.status.im_used // 226
http.status.multiple_choices // 300
http.status.moved_permanently // 301
http.status.found // 302
http.status.see_other // 303
http.status.not_modified // 304
http.status.use_proxy // 305
http.status.temporary_redirect // 307
http.status.permanent_redirect // 308
http.status.bad_request // 400
http.status.unauthorized // 401
http.status.payment_required // 402
http.status.forbidden // 403
http.status.not_found // 404
http.status.method_not_allowed // 405
http.status.not_acceptable // 406
http.status.proxy_auth_required // 407
http.status.request_timeout // 408
http.status.conflict // 409
http.status.gone // 410
http.status.length_required // 411
http.status.precondition_failed // 412
http.status.payload_too_large // 413
http.status.uri_too_long // 414
http.status.unsupported_media_type // 415
http.status.range_not_satisfiable // 416
http.status.expectation_failed // 417
http.status.teapot // 418
http.status.misdirected_request // 421
http.status.unprocessable_entity // 422
http.status.locked // 423
http.status.failed_dependency // 424
http.status.too_early // 425
http.status.upgrade_required // 426
http.status.precondition_required // 428
http.status.too_many_requests // 429
http.status.request_header_fields_too_large // 431
http.status.unavailable_for_legal_reasons // 451
http.status.internal_server_error // 500
http.status.not_implemented // 501
http.status.bad_gateway // 502
http.status.service_unavailable // 503
http.status.gateway_timeout // 504
http.status.http_version_not_supported // 505
http.status.variant_also_negotiates // 506
http.status.insufficient_storage // 507
http.status.loop_detected // 508
http.status.not_extended // 510
http.status.network_authentication_required // 511
```

## methods

`methods` is an object that provides constants for HTTP methods.

```ts
http.methods.get // "GET"
http.methods.post // "POST"
http.methods.put // "PUT"
http.methods.delete // "DELETE"
http.methods.patch // "PATCH"
http.methods.head // "HEAD"
http.methods.options // "OPTIONS"
http.methods.trace // "TRACE"
http.methods.connect // "CONNECT"
```


