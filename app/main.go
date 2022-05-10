package main

import (
	"context"
	"github.com/re-star-ru/shtrih-m-driver/app/configs"
	"github.com/re-star-ru/shtrih-m-driver/app/kkt/kktpool"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/re-star-ru/shtrih-m-driver/app/rest"
)

var version = "unknown"

func main() {
	log.Logger = log.Logger.With().
		Caller().
		Logger()
	//.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.TimeFieldFormat = time.StampMilli

	log.Info().Msgf("version: %v", version)

	// addr
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}

	// tracing
	ctx := context.Background()
	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init exporter")
	}
	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("KKTMainService")
	_, span := tracer.Start(ctx, "init service")
	span.End()
	//

	// app init

	var pool kktpool.KKTPool
	pool, err = kktpool.NewPool(configs.ConfKKT{
		"EV-S": configs.Ck{Addr: "10.51.0.71:7778", Inn: "263209745357"},
		"SM-S": configs.Ck{Addr: "10.51.0.72:7778", Inn: "262804786800"},

		"EV-N": configs.Ck{Addr: "10.51.0.73:7778", Inn: "263209745357"},
		"SM-N": configs.Ck{Addr: "10.51.0.74:7778", Inn: "262804786800"},
	})

	if err != nil {
		log.Fatal().Err(err).Send()

		return
	}

	service := rest.New(pool, addr)
	service.Run()
}

var tracer trace.Tracer

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Kkt service main"),
			semconv.ServiceVersionKey.String("0.0.1"),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
