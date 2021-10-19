package panick

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func RecoverFromPanic(ctx context.Context) {
	if r := recover(); r != nil {
		span, _ := opentracing.StartSpanFromContext(ctx, "panic")
		defer span.Finish()
		ext.Error.Set(span, true)
		span.LogFields(
			log.String("error.kind", "panic"),
			log.Object("stack", r),
		)
	}
}
