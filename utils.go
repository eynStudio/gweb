package gweb

import (
	"net/http"
)

type Handler func(c *Ctx)

type Interceptor struct {
	Before Handler
	After  Handler
}

type Resp struct {
	http.ResponseWriter
}
