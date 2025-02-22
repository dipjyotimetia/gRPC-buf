package constant

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var Tracer trace.Tracer

func init() {
	Tracer = otel.Tracer(os.Getenv("OTEL_SERVICE_NAME"))
}
