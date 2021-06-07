# Chain

Позволяет последовательно выполнять цепочку middleware, например:

```go
package main

import (
	"log"

	"google.golang.org/grpc"

	"github.com/city-mobil/gobuns/grpcext/middleware"
	"github.com/city-mobil/gobuns/grpcext/middleware/access"
	"github.com/city-mobil/gobuns/zlog"
)

func main() {
	logger, err := zlog.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer(
		middleware.WithUnaryServerChain(
			access.UnaryServerInterceptor(logger),
		),
	)

	// Register service...
}
```