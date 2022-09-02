package tracygo

import (
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

func WriteHeaders(ctx *atreugo.RequestCtx) error{
	correlationId := ctx.UserValue("X-Correlation-ID")
	requestId := ctx.UserValue("X-Request-ID")

	header := ctx.Response.Header.Add("X-Correlation-ID",correlationId)

}
