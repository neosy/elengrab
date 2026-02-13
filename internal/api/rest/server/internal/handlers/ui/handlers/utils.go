package handlers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

func getUserIDFromContext(ctx *fasthttp.RequestCtx) (uuid.UUID, error) {
	userID := ctx.UserValue(userIDKey)
	if userID == nil {
		return uuid.Nil, fmt.Errorf("user id not found")
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user id has invalid type")
	}
	return userUUID, nil
}
