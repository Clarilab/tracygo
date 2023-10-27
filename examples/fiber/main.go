package main

import (
	"fmt"
	"log"
	"sync"

	fiberTracyGo "github.com/Clarilab/tracygo/middleware/fiber/v2"
	restyTracyGo "github.com/Clarilab/tracygo/middleware/resty/v2"
	"github.com/Clarilab/tracygo/v2"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	tracer := tracygo.New(tracygo.WithCorrelationID("my-correlation-id-key"), tracygo.WithRequestID("my-request-id-key"))

	apiRestyClient := resty.New()
	apiRestyClient.OnBeforeRequest(restyTracyGo.RestyCheckTracingIDs(tracer))

	api := NewAPI(tracer, apiRestyClient)

	router := fiber.New(fiber.Config{DisableStartupMessage: true})

	router.Use(fiberTracyGo.CheckTracingIDs(tracer))
	router.Get("/hello-world", api.FiberHandler(wg))

	go func() { log.Fatal(router.Listen(":8080")) }()

	router2 := fiber.New(fiber.Config{DisableStartupMessage: true})

	router2.Use(fiberTracyGo.CheckTracingIDs(tracer))
	router2.Get("/hello-world-2", api.FiberHandler2(wg))

	go func() { log.Fatal(router2.Listen(":8081")) }()

	client := resty.New()
	client.OnBeforeRequest(restyTracyGo.RestyCheckTracingIDs(tracer))

	ctx := &atreugo.RequestCtx{
		RequestCtx: &fasthttp.RequestCtx{},
	}

	ctx.Init(&fasthttp.Request{}, nil, nil)
	ctx.SetUserValue(tracer.CorrelationIDKey(), "Zitronenbaum")

	_, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Println(err.Error())
	}

	wg.Wait()
}

// API is the api type.
type API struct {
	tracer      *tracygo.TracyGo
	restyClient *resty.Client
}

// NewAPI creates a new API.
func NewAPI(tracer *tracygo.TracyGo, restyClient *resty.Client) *API {
	return &API{
		tracer:      tracer,
		restyClient: restyClient,
	}
}

func (a *API) FiberHandler(wg *sync.WaitGroup) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		fmt.Printf("HelloWorld: X-Correlation-ID = %s\n", a.tracer.CorrelationIDromContext(ctx.Context())) // Zitronenbaum
		fmt.Printf("HelloWorld: X-Request-ID = %s\n", a.tracer.RequestIDFromContext(ctx.Context()))        // generated

		_, err := a.restyClient.R().
			SetContext(ctx.Context()).
			EnableTrace().
			Get("http://localhost:8081/hello-world-2")
		if err != nil {
			fmt.Println(err.Error())
		}

		wg.Done()

		return ctx.SendStatus(fiber.StatusOK)
	}
}

func (a *API) FiberHandler2(wg *sync.WaitGroup) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		fmt.Printf("HelloWorld2: X-Correlation-ID = %s\n", a.tracer.CorrelationIDromContext(ctx.Context())) // Zitronenbaum
		fmt.Printf("HelloWorld2: X-Request-ID = %s\n", a.tracer.RequestIDFromContext(ctx.Context()))        // new generated

		wg.Done()

		return nil
	}
}
