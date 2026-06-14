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
		exceptionx.WithMessage("Function parameter is a null pointer"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	FUNCTION_CONTAINS_EMPTY_FIELDS = exceptions.AddSeq(
		"FUNCTION_CONTAINS_EMPTY_FIELDS",
		exceptionx.WithMessage("Function parameter contains empty fields"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	QUERY_RETURNED_EMPTY_RESULT = exceptions.AddSeq(
		"QUERY_RETURNED_EMPTY_RESULT",
		exceptionx.WithMessage("A system error. The query returned an empty pointer"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)
	EMPTY_RESPONSE = exceptions.AddSeq(
		"EMPTY_RESPONSE",
		exceptionx.WithMessage("The request returned an empty response"),
		exceptionx.WithHTTPStatus(fasthttp.StatusInternalServerError),
	)

	UNAUTHORIZED = exceptions.AddNum(
		1050,
		"UNAUTHORIZED",
		exceptionx.WithMessage("Authentication is required to perform this operation"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnauthorized),
	)

	INVALID_CREDENTIALS = exceptions.AddNum(
		1051,
		"INVALID_CREDENTIALS",
		exceptionx.WithMessage("Invalid username or password"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnauthorized),
	)

	AUTH_TOKEN_MISSING = exceptions.AddNum(
		1052,
		"AUTH_TOKEN_MISSING",
		exceptionx.WithMessage("authentication token is missing"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnauthorized),
	)

	INVALID_AUTH_TOKEN = exceptions.AddNum(
		1053,
		"INVALID_AUTH_TOKEN",
		exceptionx.WithMessage("Authentication token is invalid"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnauthorized),
	)

	AUTH_TOKEN_EXPIRED = exceptions.AddNum(
		1054,
		"AUTH_TOKEN_EXPIRED",
		exceptionx.WithMessage("Authentication token has expired"),
		exceptionx.WithHTTPStatus(fasthttp.StatusUnauthorized),
	)

	INVALID_REQUEST = exceptions.AddNum(
		1100,
		"INVALID_REQUEST",
		exceptionx.WithMessage("Request contains invalid or missing data"),
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
	DOWNLOAD_NOT_FOUND = exceptions.AddSeq(
		"DOWNLOAD_NOT_FOUND",
		exceptionx.WithMessage("Media download not found"),
		exceptionx.WithHTTPStatus(fasthttp.StatusNotFound),
	)
	FILE_NOT_FOUND = exceptions.AddSeq(
		"FILE_NOT_FOUND",
		exceptionx.WithMessage("File not found"),
		exceptionx.WithHTTPStatus(fasthttp.StatusNotFound),
	)
	DOWNLOAD_ID_IS_NIL = exceptions.AddSeq(
		"DOWNLOAD_ID_IS_NIL",
		exceptionx.WithMessage("invalid download ID: nil UUID"),
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
