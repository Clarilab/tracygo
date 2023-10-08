package main

import (
	"fmt"
	"log"
	"sync"

	atreugoTracyGo "github.com/Clarilab/tracygo/middleware/atreugo/v2"
	restyTracyGo "github.com/Clarilab/tracygo/middleware/resty/v2"
	"github.com/Clarilab/tracygo/v2"

	"github.com/go-resty/resty/v2"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	tracer := tracygo.New()

	router := atreugo.New(atreugo.Config{
		Addr:   "0.0.0.0:8080",
		Logger: &Logger{},
	})

	router.UseBefore(atreugoTracyGo.AtreugoCheckTracingIDs(tracer))
	router.GET("/hello-world", AtreugoHandler(wg))
	go func() { log.Fatal(router.ListenAndServe()) }()

	router2 := atreugo.New(atreugo.Config{
		Addr:   "0.0.0.0:8081",
		Logger: &Logger{},
	})

	router2.UseBefore(atreugoTracyGo.AtreugoCheckTracingIDs(tracer))
	router2.GET("/hello-world-2", AtreugoHandler2(wg))
	go func() { log.Fatal(router2.ListenAndServe()) }()

	client := resty.New()
	client.OnBeforeRequest(restyTracyGo.RestyCheckTracingIDs(tracer))

	ctx := &atreugo.RequestCtx{
		RequestCtx: &fasthttp.RequestCtx{},
	}
	ctx.Init(&fasthttp.Request{}, nil, nil)
	ctx.SetUserValue(tracygo.CorrelationID("X-Correlation-ID"), "Zitronenbaum")

	_, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Println(err.Error())
	}

	wg.Wait()
}

func AtreugoHandler(wg *sync.WaitGroup) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {

		fmt.Printf("HelloWorld: X-Correlation-ID = %s\n", string(ctx.Request.Header.Peek("X-Correlation-ID"))) // Zitronenbaum
		fmt.Printf("HelloWorld: X-Request-ID = %s\n", string(ctx.Request.Header.Peek("X-Request-ID")))         // generated

		tracy := tracygo.New()
		client := resty.New()
		client.OnBeforeRequest(restyTracyGo.RestyCheckTracingIDs(tracy))

		_, err := client.R().
			SetContext(ctx).
			EnableTrace().
			Get("http://localhost:8081/hello-world-2")
		if err != nil {
			fmt.Println(err.Error())
		}

		wg.Done()

		return ctx.JSONResponse(nil, 200)
	}
}

func AtreugoHandler2(wg *sync.WaitGroup) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		fmt.Printf("HelloWorld2: X-Correlation-ID = %s\n", string(ctx.Request.Header.Peek("X-Correlation-ID"))) // Zitronenbaum
		fmt.Printf("HelloWorld2: X-Request-ID = %s\n", string(ctx.Request.Header.Peek("X-Request-ID")))         // generated

		wg.Done()

		return nil
	}
}
