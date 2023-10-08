# TracyGo

## INFO

TracyGo provides functionality for tracing a correlation identifier through multiple go microservices.

The library includes middlewares for the following frameworks:
- [Atreugo](https://github.com/savsgio/atreugo)
- [Fiber](https://github.com/gofiber/fiber)
- [Resty](https://github.com/go-resty/resty)

### Supported Go Versions

This library supports the most recent Go, currently  **1.21**.

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
