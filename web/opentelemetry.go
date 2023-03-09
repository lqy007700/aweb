package web

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func (m *MiddlewareBuilder) Builder() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			spanCtx, span := m.tracer.Start(c.R.Context(), "my-span")
			defer span.End()
			c.R = c.R.WithContext(spanCtx)
			span.SetAttributes(attribute.String("http.method", c.R.Method))
			span.SetAttributes(attribute.String("http.path", c.R.URL.Path[:256]))
			next(c)
		}
	}
}
