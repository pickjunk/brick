# brick
[![GoDoc](https://godoc.org/github.com/pickjunk/brick?status.svg)](https://godoc.org/github.com/pickjunk/brick)

A Simple Golang HTTP Framework.

### Features

- High-performance router, base on [httprouter](https://github.com/julienschmidt/httprouter)
- Pretty middewares
- Uniform logger, base on [zerolog](https://github.com/rs/zerolog)
- Opentracing, integrate jaeger-client
- Graphql, base on [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go)
- CORS, base on [rs/cors](https://github.com/rs/cors)

### Quick Start

```golang
// main.go
package main

import (
  "net/http"

  dbr "github.com/gocraft/dbr"
  httprouter "github.com/julienschmidt/httprouter"
  b "github.com/pickjunk/brick"
  bd "github.com/pickjunk/brick/dbr"
)

func main() {
  r := b.New()

  r.GET("/:name", func(ctx context.Context) {
    w := b.Response(ctx) // http.ResponseWriter
    name := b.Param(ctx, "name") // httprouter.Params.ByName("name")

    w.Write([]byte("hello "+name+"!"))
  })

  // comment out these code if you don't use dbr
  rWithDbr := r.Prefix("/dbr").Middlewares(bd.Middleware(nil))

  rWithDbr.GET("/empty", func(ctx context.Context) {
    w := b.Response(ctx) // http.ResponseWriter
    db := bd.Dbr(ctx) // *dbr.Session

    var test struct{}
    err := db.Select(`"empty"`).LoadOneContext(ctx, &test)
    if err != nil {
      brick.log.Panic(err)
    }

    w.Write([]byte(`dbr: SELECT "empty"`))
  })

  r.ListenAndServe()
}
```

### Config With Environments

* PORT - the port listened to, optional, default 8080
* MYSQL_DSN - mysql dsn used by dbr, required if use brick/dbr
* MYSQL_MAX_IDLE - the count of mysql max idle connections, optional, default 1
* MYSQL_MAX_OPEN - the count of mysql max open connections, optional, default 1
* MAIL_HOST - the host of mail server, required if use brick/utils/mail
* MAIL_PORT - the port of mail server, required if use brick/utils/mail
* MAIL_USER - the user of mail server, required if use brick/utils/mail
* MAIL_PASSWD - the passwd of mail server, required if use brick/utils/mail

### Middlewares

```golang
r.Middlewares(
  func(ctx context.Context, next brick.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context, next brick.Handle) {
    // do something

    next(ctx)
  },
).GET(
  "/"
  func(ctx context.Context, next brick.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context) {
    // do something
  },
)
```

### Context

```golang
r.Middlewares(
  func(ctx context.Context, next brick.Handle) {
    // context with a key & value
    next(b.WithValue(ctx, "key", value))
  },
).GET(
  "/"
  func(ctx context.Context) {
    // get a value from context by a key
    value := b.Value(ctx, "key").(type) // "type" means a type assertion here
  },
)
```

### HTTP

```golang
r.GET("/:name", func(ctx context.Context) {
  r := b.Request(ctx) // *http.Request
  w := b.Response(ctx) // http.ResponseWriter
  ps := b.Params(ctx) // httprouter.Params
  value := b.Param(ctx, "key") // httprouter.Params.ByName("key")
})
```

### SubRoute (Prefix + Middlewares)

```golang
subRoute1 := r.Prefix("/sub1")

subRoute2 := r.Prefix("/sub1").Middlewares(
  func(ctx context.Context, next brick.Handle) {
    // do something

    next(ctx)
  },
  func(ctx context.Context, next brick.Handle) {
    // do something

    next(ctx)
  },
)
```

### Graphql

```golang
type resolver struct{}

var g = b.NewGraphql(&resolver{})

func init() {
  g.Schema(`
  schema {
    query: Query
  }

  type Query {
    greeting(name: String): String!
  }
  `)
}

func (r *resolver) Greeting(
  ctx context.Context,
  args struct {
    Name  *string
  },
) string {
  if args.Name == nil {
    name := "world"
    args.Name = &name
  }
  return "hello " + args.Name
}

func main() {
  r := b.New()

  r.Graphql("/graphql", g)
}
```

### Opentracing (jaeger-client)

```golang
package main

import (
  b "github.com/pickjunk/brick"
  config "github.com/uber/jaeger-client-go/config"
)

func main() {
  r := b.New()

  closer := b.Jaeger(&config.Configuration{
    ServiceName: "brick-example",
    Sampler: &config.SamplerConfig{
      Type:  "const",
      Param: 1,
    },
    Reporter: &config.ReporterConfig{
      LogSpans: true,
    },
  })
  defer closer.Close()

  r.ListenAndServe()
}
```

### CORS

```golang
package main

import (
  b "github.com/pickjunk/brick"
  "github.com/rs/cors"
)

func main() {
  r := b.New()

  r.CORS(cors.AllowAll()).GET("/", func(ctx context.Context) {
    // do something
  })

  r.ListenAndServe()
}
```

### Business Error

```golang
package main

import (
  b "github.com/pickjunk/brick"
  be "github.com/pickjunk/brick/error"
)

func main() {
  r := b.New()

  r.GET("/", func(ctx context.Context) {
    // Throw will trigger a panic, which internal recover middleware
    // will catch and unmarshal to the response content as
    // `{"code":10001,"msg":"passwd error"}`
    be.Throw(10001, "passwd error")
  })

  r.ListenAndServe()
}
```

### Logger

```golang
// create a file named logger.go in your package
// (for a convention, you should create this file in every package,
// so you can simplely call log.xxx() in every package)
package main

import (
  bl "github.com/pickjunk/brick/log"
)

var log = bl.New("component's name, like a namespace, for example, brick.main")

// main.go with the logger.go above in the same package
package main

import (
  b "github.com/pickjunk/brick"
)

func main() {
  r := b.New()

  r.GET("/", func(ctx context.Context) {
    log.Info().Msg("request access")
    // yes, the logger based on zerolog
    // more documents are [here](https://github.com/rs/zerolog)
  })

  r.ListenAndServe()
}
```
