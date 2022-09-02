package tracygo

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/savsgio/atreugo/v11"
)

func CheckRequestID(ctx *atreugo.RequestCtx) error {

	correlationId := string(ctx.Request.Header.Peek("X-Correlation-ID"))

	if correlationId == "" {
		correlationId = uuid.New().String()
	}

	requestId := uuid.New().String()

	ctx.SetUserValue("X-Correlation-ID", correlationId)
	ctx.SetUserValue("X-Request-ID", requestId)

	return ctx.Next()
}

func WriteHeader(ctx *atreugo.RequestCtx) error {
	correlationId := ctx.UserValue("X-Correlation-ID").(string)
	requestId := ctx.UserValue("X-Request-ID").(string)

	ctx.Response.Header.Set("X-Correlation-ID", correlationId)
	ctx.Response.Header.Set("X-Request-ID", requestId)

	fmt.Println(ctx.Response.Header.String())

	return ctx.Next()
}

func CheckRequestIDResty(client *resty.Client, request *resty.Request) error {
	correlationId, ok := request.Context().Value("X-Correlation-ID").(string)
	if !ok {
		return nil
	}

	requestId, ok := request.Context().Value("X-Request-ID").(string)
	if !ok {
		return nil
	}

	if correlationId != "" || requestId != "" {
		request.Header.Set("X-Correlation-ID", correlationId)
		request.Header.Set("X-Request-ID", requestId)
	}

	return nil
}
