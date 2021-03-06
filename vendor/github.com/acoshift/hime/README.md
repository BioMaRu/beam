# Hime

[![Build Status](https://travis-ci.org/acoshift/hime.svg?branch=master)](https://travis-ci.org/acoshift/hime)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/hime/badge.svg?branch=master)](https://coveralls.io/github/acoshift/hime?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/hime)](https://goreportcard.com/report/github.com/acoshift/hime)
[![GoDoc](https://godoc.org/github.com/acoshift/hime?status.svg)](https://godoc.org/github.com/acoshift/hime)

Hime is a Go Web Framework.

## Why Framework

I like net/http but... there are many duplicated code when working on multiple projects,
plus no standard. Framework creates a standard for developers.

### Why Another Framework

There're many Go framework out there. But I want a framework that works with any net/http compatible libraries seamlessly.

For example, you can choose any router, any middlewares, or handlers that work with standard library.

Other framework don't allow this. They have built-in router, framework-specific middlewares.

## Core focus

- Add standard to code
- Compatible with any net/http compatible router
- Compatible with http.Handler without code change
- Compatible with net/http middlewares without code change
- Use standard html/template for view
- Built-in core functions for build web server
- Reduce developer bugs

## What is this framework DO NOT focus

- Speed
- One framework do everything

## Example

```go
func main() {
    app := hime.New()

    app.Template().
        Parse("index", "index.tmpl", "_layout.tmpl").
        Minify()

    app.
        Routes(hime.Routes{
            "index": "/",
        }).
        BeforeRender(addHeaderRender).
        Handler(router()).
        GracefulShutdown().
        Address(":8080").
        ListenAndServe()
}

func router() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("/", hime.Handler(indexHandler))
    return middleware.Chain(
        logRequestMethod,
        logRequestURI,
    )(mux)
}

func logRequestURI(h http.Handler) http.Handler {
    return hime.Handler(func(ctx *hime.Context) error {
        log.Println(ctx.Request().RequestURI)
        return h
    })
}

func logRequestMethod(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method)
        h.ServeHTTP(w, r)
    })
}

func addHeaderRender(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        h.ServeHTTP(w, r)
    })
}

func indexHandler(ctx *hime.Context) error {
    if ctx.Request().URL.Path != "/" {
        return ctx.RedirectTo("index")
    }
    return ctx.View("index", map[string]interface{}{
        "Name": "Acoshift",
    })
}
```

### More Examples

- [Todo App](https://github.com/acoshift/todo-hime)

## Compatibility with net/http

### Handler

Hime doesn't have built-in router, you can use any http.Handler.

`hime.Handler` already implements `http.Handler`, so you can use hime's handler anywhere in your router that support http.Handler.

### Middleware

Hime use native func(http.Handler) http.Handler for middleware.
You can use any middleware that compatible with this type.

```go
func logRequestMethod(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.Method)
        h.ServeHTTP(w, r)
    })
}
```

You can also use hime's handler with middleware

```go
func logRequestURI(h http.Handler) http.Handler {
    return hime.Handler(func(ctx *hime.Context) error {
        log.Println(ctx.Request().RequestURI)
        return ctx.Handle(h)
    })
}
```

Inject data to context

```go
func injectData(h http.Handler) http.Handler {
    return hime.Handler(func(ctx *hime.Context) error {
        ctx.WithValue(ctxKeyData{}, "injected data!")
        return ctx.Handle(h)
    })
}
```

Retrieve data from context like context.Context

```go
ctx.Value(ctxKeyData{}).(string)
```

### Use hime.App as Handler

If you don't want to use hime's built-in graceful shutdown,
you can use hime.App as normal handler.

```go
func main() {
    app := hime.New()

    app.Template().
        Parse("index", "index.tmpl", "_layout.tmpl").
        Minify()

    app.
        Routes(hime.Routes{
            "index": "/",
        }).
        BeforeRender(addHeaderRender).
        Handler(router())

    http.ListenAndServe(":8080", app)
}
```

## Panic when return error

When return error from handler, hime will panic,
this mean you should handle error and only return unrecovery error.

```go
func signInHandler(ctx *hime.Context) error {
    username := r.FormValue("username")

    user, err := findUser(username)
    if err == ErrUserNotFound {
        return ctx.Status(http.StatusBadRequest).Error("invalid username")
    }
    if err != nil {
        return err // database connection error ?
    }
    ...
}
```

## Why some functions use panic

Hime try to reduce developer errors,
some error can detect while development.
Hime will panic for that type of errors.

## Init App with YAML

Hime can init using YAML.

```yaml
// app.yaml
globals:
  data1: test
routes:
  index: /
  about: /about
templates:
- dir: view
  root: layout
  delims: ["{{", "}}"]
  minify: true
  components:
  - comp/comp1.tmpl
  - comp/comp2.tmpl
  list:
    main.tmpl:
    - main.tmpl
    - _layout.tmpl
    about.tmpl: [about.tmpl, _layout.tmpl]
server:
  addr: :8080
  readTimeout: 10s
  readHeaderTimeout: 5s
  writeTimeout: 5s
  idleTimeout: 30s
  gracefulShutdown:
    timeout: 1m
    wait: 5s
```

```go
app := hime.New().
    ParseConfigFile("app.yaml")
```

### Multiple Configs

```yaml
// routes.yaml
routes:
  index: /
  about: /about
```

```yaml
// server.yaml
server:
  addr: :8080
  readTimeout: 10s
  readHeaderTimeout: 5s
  writeTimeout: 5s
  idleTimeout: 30s
  gracefulShutdown:
    timeout: 1m
    wait: 5s
```

```yaml
// template.web.yaml
dir: view
root: layout
delims: ["{{", "}}"]
minify: true
components:
- comp/comp1.tmpl
- comp/comp2.tmpl
list:
main.tmpl:
- main.tmpl
- _layout.tmpl
about.tmpl: [about.tmpl, _layout.tmpl]
```

```go
app := hime.New().
    ParseConfigFile("routes.yaml").
    ParseConfigFile("server.yaml")

app.Template().ParseConfigFile("template.web.yaml")
```

## Graceful Shutdown Multiple Apps

Hime can handle graceful shutdown for multiple Apps.

```go
app1 := hime.New().ParseConfigFile("app1.yaml")
app2 := hime.New().ParseConfigFile("app2.yaml")

probe := probehandler.New() // github.com/acoshift/probehandler
health := http.NewServeMux()
health.Handle("/readiness", probe)
health.Handle("/liveness", probehandler.Success())
go http.ListenAndServe(":18080", health)

err := hime.Merge(app1, app2).
    GracefulShutdown().
    Notify(probe.Fail).
    ListenAndServe()
if err != nil {
    log.Fatal(err)
}
```

## Useful handlers and middlewares

- [acoshift/middleware](https://github.com/acoshift/middleware)
- [acoshift/header](https://github.com/acoshift/header)
- [acoshift/session](https://github.com/acoshift/session)
- [acoshift/flash](https://github.com/acoshift/flash)
- [acoshift/webstatic](https://github.com/acoshift/webstatic)
- [acoshift/httprouter](https://github.com/acoshift/httprouter)
- [acoshift/redirecthttps](https://github.com/acoshift/redirecthttps)
- [acoshift/cachestatic](https://github.com/acoshift/cachestatic)
- [acoshift/methodmux](https://github.com/acoshift/methodmux)

## FAQ

- Will hime support automated let's encrypt ?
  > No, you can use hime as normal handler and use other lib to do this.

- Do you use hime in production ?
  > Yes, why not ? :D

## License

MIT
