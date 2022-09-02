package main

import (
	"fmt"
	"dev.azure.com/kycnow/kycnow/_git/tracygo"
	"github.com/savsgio/atreugo/v11"

)

func main() {

	router := atreugo.New(atreugo.Config{
		Addr : "0.0.0.0:8080",
	})

	router.UseBefore(tracygo.CheckRequestID)
	router.UseAfter(tracygo.WriteHeaders)
	router.GET("/hello-world", SomeFunction)
	router.ListenAndServe()
	
	//client := resty.New()
	
	//_, err := client.R().
	//	EnableTrace().
	//	Get("http://localhost:8080")
	//if err != nil {
	//	fmt.Printf("error")
	//}

	//fmt.Printf(string(resp))

}

func SomeFunction(ctx *atreugo.RequestCtx) error {
	fmt.Println("Hello world")
	fmt.Println(ctx.UserValue("X-Request-ID").(string))
	fmt.Println(ctx.UserValue("X-Correlation-ID").(string))
	return nil
}