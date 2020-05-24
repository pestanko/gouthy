package shared

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

func NewOperationContext() context.Context {
	opId := uuid.NewV4()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "operation_id", opId.String())
	return ctx
}
