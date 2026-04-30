package exceptions

import (
	"github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"
	"github.com/valyala/fasthttp"
)

var (
	exceptions = exceptionx.NewExceptions(20)

	FUNCTION_PARAMETER_NULL_POINTER = exceptions.AddNum(
		1000,
		"FUNCTION_PARAMETER_NULL_POINTER",
		exceptionx.WithMessage("function parameter is a null pointer"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	FUNCTION_CONTAINS_EMPTY_FIELDS = exceptions.AddSeq(
		"FUNCTION_CONTAINS_EMPTY_FIELDS",
		exceptionx.WithMessage("function parameter contains empty fields"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	QUERY_RETURNED_EMPTY_RESULT = exceptions.AddSeq(
		"QUERY_RETURNED_EMPTY_RESULT",
		exceptionx.WithMessage("a system error. The query returned an empty pointer"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	EMPTY_RESPONSE = exceptions.AddSeq(
		"EMPTY_RESPONSE",
		exceptionx.WithMessage("the request returned an empty response"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)

	INVALID_REQUEST = exceptions.AddNum(
		1050,
		"INVALID_REQUEST",
		exceptionx.WithMessage("request contains invalid or missing data"),
		exceptionx.WithHTTPStatus(fasthttp.StatusBadRequest),
	)

	DEMO_MODE_RESTRICTION = exceptions.AddNum(
		2000,
		"DEMO_MODE_RESTRICTION",
		exceptionx.WithMessage("Operation not allowed in demo mode"),
		exceptionx.WithHTTPStatus(fasthttp.StatusForbidden),
	)

	DOWNLOADER_EMPTY_RESPONSE = exceptions.AddNum(
		2100,
		"DOWNLOADER_EMPTY_RESPONSE",
		exceptionx.WithMessage("Service downloader returned an empty value"),
		exceptionx.WithHTTPStatus(fasthttp.StatusBadGateway),
	)
	FILE_NOT_FOUND = exceptions.AddSeq(
		"FILE_NOT_FOUND",
		exceptionx.WithMessage("File not found"),
		exceptionx.WithHTTPStatus(fasthttp.StatusNotFound),
	)
	FILE_ID_IS_NIL = exceptions.AddSeq(
		"FILE_ID_IS_NIL",
		exceptionx.WithMessage("invalid file ID: nil UUID"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnprocessableEntity),
	)
	QUEUE_PUBLISH_FAILED = exceptions.AddSeq(
		"QUEUE_PUBLISH_FAILED",
		exceptionx.WithMessage("failed to enqueue download task"),
		exceptionx.WithHTTPStatus(fasthttp.StatusServiceUnavailable),
	)
	THUMBNAIL_NOT_FOUND = exceptions.AddSeq(
		"THUMBNAIL_NOT_FOUND",
		exceptionx.WithMessage("Thumbnail not found"),
		exceptionx.WithHTTPStatus(fasthttp.StatusNotFound),
	)
)
