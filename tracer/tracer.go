package tracer

import (
	"context"
	"crypto/tls"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type TransportProtocolType string

const (
	TransportProtocolHTTP TransportProtocolType = "HTTP"
	TransportProtocolGRPC TransportProtocolType = "GRPC"
)

type TraceConf struct {
	ServiceName       string                `mapstructure:"serviceName"`
	CollectorURL      string                `mapstructure:"collectorURL"`
	Insecure          bool                  `mapstructure:"insecure"`
	TransportProtocol TransportProtocolType `mapstructure:"transportProtocol"`
}

func InitTracer(conf TraceConf) func(context.Context) error {
	var (
		serviceName  = conf.ServiceName
		collectorURL = conf.CollectorURL
		insecure     = conf.Insecure
	)
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("[InitTracer] recover failed")
		}
	}()

	var exporter *otlptrace.Exporter
	var err error
	switch conf.TransportProtocol {
	case TransportProtocolHTTP:
		secureOption := otlptracehttp.WithTLSClientConfig(&tls.Config{})
		if insecure {
			secureOption = otlptracehttp.WithInsecure()
		}
		exporter, err = otlptrace.New(
			context.Background(),
			otlptracehttp.NewClient(
				otlptracehttp.WithEndpointURL(collectorURL),
				secureOption,
			),
		)
	default:
		secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
		if insecure {
			secureOption = otlptracegrpc.WithInsecure()
		}
		exporter, err = otlptrace.New(
			context.Background(),
			otlptracegrpc.NewClient(
				secureOption,
				otlptracegrpc.WithEndpoint(collectorURL),
			),
		)

	}
	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Fatal("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resources),
		),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return exporter.Shutdown
}

func TraceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func InitContextWithTrace(ctx context.Context, t32, s16 string) context.Context {
	tid, err := trace.TraceIDFromHex(t32)
	if err != nil {
		tid, _ = trace.TraceIDFromHex("01000000000000000000000000000000")
	}
	sid, err := trace.SpanIDFromHex(s16)
	if err != nil {
		sid, _ = trace.SpanIDFromHex("0200000000000000")
	}

	sc := trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid, SpanID: sid})
	return trace.ContextWithSpanContext(ctx, sc)
}
