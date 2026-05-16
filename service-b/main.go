package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	"github.com/teichx/go-expert-weather/internal/handlers"
	"github.com/teichx/go-expert-weather/pkg/viacepclient"
	"github.com/teichx/go-expert-weather/pkg/weatherapiclient"
)

func initTracer(ctx context.Context) (*trace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("service-b"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return tp, nil
}

func main() {
	_ = godotenv.Load()
	weatherAPIKey, ok := os.LookupEnv("WEATHER_API_KEY")
	if !ok {
		log.Fatal("WEATHER_API_KEY environment variable not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer tp.Shutdown(ctx)

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	viaCEPClient := viacepclient.New(httpClient)
	weatherClient := weatherapiclient.New(weatherAPIKey, httpClient)

	h := handlers.New(viaCEPClient, weatherClient)

	handler := otelhttp.NewHandler(http.HandlerFunc(h.PostWeather), "POST /")
	http.Handle("/", handler)

	addr := ":8081"
	log.Printf("Service B listening on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
