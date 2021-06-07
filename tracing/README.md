# Tracing

### Config example:

```yaml
jeager:
    service_name: "foo"
    disabled: false         # default
    agent: "localhost:6831" # default
```

### Init OT

```go
import (
    "log"

    "github.com/city-mobil/gobuns/config"
    "github.com/city-mobil/gobuns/tracing"
    "github.com/city-mobil/gobuns/tracing/jaegercfg"
)

tracerConfFn := jaegercfg.JaegerConfig("")

err := config.InitOnce()
if err != nil {
    log.Fatal("failed init config")
}

closer, err := tracing.InitGlobalTracerFromConfig(tracerConfFn())
if err != nil {
    log.Fatal("failed to init Open Tracing ")
}
defer closer.Close()
```

### MySQL

```go
import (
    "database/sql"

    // register "traceable-mysql" driver
    _ "github.com/city-mobil/gobuns/tracing/mysql"
)

db, err := sql.Open("traceable-mysql", "user:password@/dbname")
```

### HTTP middleware

```go
import (
    "github.com/opentracing/opentracing-go"
    "github.com/city-mobil/gobuns/tracing"
)
tracer := opentracing.GlobalTracer()
mw := tracing.NewHTTPMiddleware(tracer)


handlerFunc := func(w http.ResponseWriter, r *http.Request){}
http.HandleFunc("/handler1", mw.HandlerFunc(handlerFunc))


type handler struct{}
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
http.Handle("/handler2", mw.Handler(&handler{}))
```
