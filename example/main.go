package main

import (
	"fmt"

	"github.com/Clarilab/tracygo"
	"github.com/go-resty/resty/v2"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
)

func main() {
	tracy := tracygo.New()

	router := atreugo.New(atreugo.Config{
		Addr: "0.0.0.0:8080",
	})

	router.UseBefore(tracy.AtreugoCheckTracingIDs)
	router.GET("/hello-world", SomeFunction)
	go router.ListenAndServe()

	router2 := atreugo.New(atreugo.Config{
		Addr: "0.0.0.0:8081",
	})

	router2.UseBefore(tracy.AtreugoCheckTracingIDs)
	router2.GET("/hello-world-2", SomeFunction2)
	go router2.ListenAndServe()

	client := resty.New()
	client.OnBeforeRequest(tracy.RestyCheckTracingIDs)

	ctx := &atreugo.RequestCtx{
		RequestCtx: &fasthttp.RequestCtx{},
	}
	ctx.Init(&fasthttp.Request{}, nil, nil)
	ctx.SetUserValue("X-Correlation-ID", "Zitronenbaum")

	_, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Println(err.Error())
	}

	select {}
}

func SomeFunction(ctx *atreugo.RequestCtx) error {
	fmt.Printf("HelloWorld: X-Correlation-ID = %s\n", string(ctx.Request.Header.Peek("X-Correlation-ID"))) // Zitronenbaum
	fmt.Printf("HelloWorld: X-Request-ID = %s\n", string(ctx.Request.Header.Peek("X-Request-ID")))         // generated

	tracy := tracygo.New()
	client := resty.New()
	client.OnBeforeRequest(tracy.RestyCheckTracingIDs)

	_, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8081/hello-world-2")
	if err != nil {
		fmt.Println(err.Error())
	}

	return ctx.JSONResponse(nil, 500)
}

func SomeFunction2(ctx *atreugo.RequestCtx) error {
	fmt.Printf("HelloWorld2: X-Correlation-ID = %s\n", string(ctx.Request.Header.Peek("X-Correlation-ID"))) // Zitronenbaum
	fmt.Printf("HelloWorld2: X-Request-ID = %s\n", string(ctx.Request.Header.Peek("X-Request-ID")))         // generated

	return nil
}
