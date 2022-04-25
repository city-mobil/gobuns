package handlers

import (
	"context"
	"net/http"
)

type Context interface {
	context.Context
	HTTPRequest() *http.Request
}

type baseContext struct {
	context.Context
	httpReq *http.Request
}

func (c *baseContext) HTTPRequest() *http.Request {
	return c.httpReq
}
