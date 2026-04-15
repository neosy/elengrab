package exceptions

import (
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/valyala/fasthttp"
)

var (
	exceptions = exceptionx.NewExceptions(10)

	DEMO_MODE_RESTRICTION = exceptions.AddNum(
		1101,
		"DEMO_MODE_RESTRICTION",
		exceptionx.WithMessage("Operation not allowed in demo mode"),
		exceptionx.WithHTTPStatus(fasthttp.StatusForbidden),
	)
)
