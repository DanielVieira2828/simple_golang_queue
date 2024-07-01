package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/savi2w/simple-queue/rmq"
)

type Router struct {
	Broker *rmq.Broker
	Server *echo.Echo
}

func (r *Router) Handler(ctx echo.Context) error {

	response, err := r.Broker.MakeRequest()
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.String(http.StatusOK, response.Uuid)
}
