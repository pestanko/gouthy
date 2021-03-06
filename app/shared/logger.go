package shared

import (
	"context"
	log "github.com/sirupsen/logrus"
)

func GetLogger(ctx context.Context) *log.Entry {
	var entry = log.WithFields(log.Fields{})
	//log.SetFormatter(&log.JSONFormatter{})

	if opId := ctx.Value("operation_id"); opId != nil {
		entry = entry.WithField("operation_id", opId.(string))
	}
	id := ctx.Value("identity")
	if id != nil {
		entry = entry.WithFields(id.(LoggingIdentity).LogFields())
	}
	return entry
}


type LoggingIdentity interface {
	LogFields() log.Fields
}
