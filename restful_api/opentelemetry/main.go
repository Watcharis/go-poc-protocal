package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

// initTracer initializes the OpenTelemetry Tracer Provider.
func initTracer() (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("golang-microservice"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

// callAnotherService demonstrates making an HTTP request with context propagation.
func callAnotherService(ctx context.Context) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	req, _ := http.NewRequestWithContext(ctx, "GET", "http://localhost:8777/health", nil)
	_, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
	}
}

// mainHandler handles incoming HTTP requests.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Call another service with propagated context
	callAnotherService(ctx)

	w.Write([]byte("Request complete!"))
}

func main() {
	// Initialize the Tracer Provider
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() { _ = tp.Shutdown(context.Background()) }()

	// Wrap the main handler with OpenTelemetry middleware
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(mainHandler), "main-handler"))

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
