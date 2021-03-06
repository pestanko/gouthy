package shared

import (
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

func NewOperationContext() context.Context {
	opId := uuid.NewV4()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "operation_id", opId.String())
	return ctx
}

func NewContextWithConfiguration(config *AppConfig) context.Context {
	ctx := NewOperationContext()
	ctx = context.WithValue(ctx, "config", config)
	return ctx
}

func ToJSON(obj interface{}) string {
	data, _ := json.Marshal(obj)
	return string(data)
}

func ToJSONIndent(obj interface{}) string {
	data, _ := json.MarshalIndent(obj, "", "  ")
	return string(data)
}

func GetConfig(ctx context.Context) *AppConfig {
	value := ctx.Value("config")
	if value == nil {
		return nil
	}
	return value.(*AppConfig)
}
