# TracyGo

## INFO

TracyGo provides functionality for tracing a correlation identifier through multiple go microservices.

The library includes middlewares for the following frameworks:
- [Atreugo](https://github.com/savsgio/atreugo)
- [Echo](https://github.com/labstack/echo)
- [Fiber](https://github.com/gofiber/fiber)
- [Resty](https://github.com/go-resty/resty)
- [net/http](https://pkg.go.dev/net/http)
- [EventHorizon](https://github.com/looplab/eventhorizon)

### Supported Go Versions

This library supports the most recent Go, currently  **1.22.3**.

## How to install

```bash
go get github.com/Clarilab/tracygo/v2
```

## How to import
```go
import "github.com/Clarilab/tracygo/v2"
```

## How to use
In the examples folder you will find some example applications.
One using the atreugo middleware, one using the fiber middleware.
Both examples also demonstrate how to use the resty middleware.
