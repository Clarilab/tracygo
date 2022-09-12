package main

import (
	"context"
	"fmt"

	"dev.azure.com/kycnow/kycnow/_git/tracygo"
	"github.com/go-resty/resty/v2"
	"github.com/savsgio/atreugo/v11"
)

func main() {
	tracy := tracygo.New()

	router := atreugo.New(atreugo.Config{
		Addr: "0.0.0.0:8080",
	})

	router.UseBefore(tracy.CheckRequestID)
	router.UseAfter(tracy.WriteHeader)
	router.GET("/hello-world", SomeFunction)
	go router.ListenAndServe()

	client := resty.New()
	client.OnBeforeRequest(tracy.CheckTracingIDs)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "X-Correlation-ID", "Zitronenbaum")
	ctx = context.WithValue(ctx, "X-Request-ID", "Dies das")
	_, err := client.R().
		SetContext(ctx).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SomeFunction(ctx *atreugo.RequestCtx) error {
	fmt.Println("Hello world")
	return nil
}
