package logging

import (
	"net/http"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LoggerWithRequest(r *http.Request) otelzap.LoggerWithCtx {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, _ := config.Build()
	span := trace.SpanContextFromContext(r.Context())
	logger = logger.With(
		zap.String("trace_id", span.TraceID().String()),
		zap.String("span_id", span.SpanID().String()),
		zap.String("request_uri", r.RequestURI),
		zap.String("method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("url", r.URL.String()),
		zap.String("referer", r.Header.Get("Referer")),
		zap.String("user-agent", r.Header.Get("User-Agent")),
		zap.String("x-request-id", r.Header.Get("X-Request-ID")),
	)
	return otelzap.New(logger).Ctx(r.Context())
}

func LoggerWithSpanContext(span trace.SpanContext) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	logger, _ := config.Build()
	return logger.With(
		zap.String("trace_id", span.TraceID().String()),
		zap.String("span_id", span.SpanID().String()),
	)
}
