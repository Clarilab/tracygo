package main

import (
	"context"
	"fmt"

	"dev.azure.com/kycnow/kycnow/_git/tracygo"
	"github.com/go-resty/resty/v2"
	"github.com/savsgio/atreugo/v11"
)

func main() {

	//router := atreugo.New(atreugo.Config{
	//	Addr: "0.0.0.0:8080",
	//})

	// router.UseBefore(tracygo.CheckRequestID)
	// router.UseAfter(tracygo.WriteHeader)
	// router.GET("/hello-world", SomeFunction)
	// router.ListenAndServe()

	client := resty.New()
	client.OnBeforeRequest(tracygo.CheckRequestIDResty)
	_, err := client.R().
		SetContext(context.WithValue(context.Background(), "X-Correlation-ID", "Zitronenbaum")).
		EnableTrace().
		Get("http://localhost:8080/hello-world")
	if err != nil {
		fmt.Printf("error")
	}

}

func SomeFunction(ctx *atreugo.RequestCtx) error {
	fmt.Println("Hello world")
	return nil
}
