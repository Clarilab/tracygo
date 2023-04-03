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

	client := resty.New()
	client.OnBeforeRequest(tracy.RestyCheckTracingIDs)

	ctx := &atreugo.RequestCtx{
		RequestCtx: &fasthttp.RequestCtx{},
	}
	ctx.Init(&fasthttp.Request{}, nil, nil)
	ctx.SetUserValue("X-Correlation-ID", "Zitronenbaum")
	ctx.SetUserValue("X-Request-ID", "1234567890")

	resp, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(resp.Header().Get("X-Correlation-ID"))
	fmt.Println(resp.Header().Get("X-Request-ID"))

	select {}
}

func SomeFunction(ctx *atreugo.RequestCtx) error {
	fmt.Println("Hello world")
	return nil
}
