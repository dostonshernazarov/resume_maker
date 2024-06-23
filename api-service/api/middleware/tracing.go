package middleware

// import (
// 	"bufio"
// 	"errors"
// 	"net"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"go.opencensus.io/trace"
// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"

// 	// "github.com/dostonshernazarov/resume_maker/api/response"
// 	// "github.com/dostonshernazarov/resume_maker/internal/pkg/otlp"
// )

// // Tracing middleware function
// func Tracing(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	tracer := otel.GetTracerProvider().Tracer("api-service-booking")
// 	ctx, span := tracer.Start(ctx, "HTTP "+c.Request.Method+" "+c.FullPath(), trace.WithSpanKind(trace.SpanKindServer))
// 	defer span.End()

// 	c.Writer.Header().Add("TraceID", span.SpanContext().TraceID().String())

// 	rw := &responseWriter{c.Writer, http.StatusOK}

// 	c.Writer = rw
// 	c.Request = c.Request.WithContext(ctx)

// 	c.Next()

// 	span.SetAttributes(
// 		attribute.String("http.method", c.Request.Method),
// 		attribute.String("http.url", c.FullPath()),
// 		attribute.Int("http.status_code", rw.statusCode),
// 	)
// }

// type responseWriter struct {
// 	gin.ResponseWriter
// 	statusCode int
// }

// func (rw *responseWriter) WriteHeader(statusCode int) {
// 	rw.statusCode = statusCode
// 	rw.ResponseWriter.WriteHeader(statusCode)
// }

// func (rw *responseWriter) Write(data []byte) (int, error) {
// 	return rw.ResponseWriter.Write(data)
// }

// func (rw *responseWriter) WriteString(s string) (int, error) {
// 	return rw.ResponseWriter.WriteString(s)
// }

// func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
// 	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
// 		return hijacker.Hijack()
// 	}
// 	return nil, nil, errors.New("response writer does not support hijacking")
// }

// func (rw *responseWriter) CloseNotify() <-chan bool {
// 	if notifier, ok := rw.ResponseWriter.(http.CloseNotifier); ok {
// 		return notifier.CloseNotify()
// 	}
// 	return make(chan bool)
// }
