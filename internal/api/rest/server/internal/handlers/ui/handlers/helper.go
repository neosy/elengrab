package handlers

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type requestFilters map[string]string

func parseFilters(ctx *fasthttp.RequestCtx) requestFilters {
	filters := make(requestFilters)

	for key, value := range ctx.QueryArgs().All() {
		k := string(key)
		v := string(value)

		prefix := "filter["
		suffix := "]"

		if !strings.HasPrefix(k, prefix) || !strings.HasSuffix(k, suffix) {
			continue
		}

		// filter[name] → name
		field := k[len(prefix) : len(k)-len(suffix)]
		if field == "" {
			continue
		}

		filters[field] = v
	}

	return filters
}
