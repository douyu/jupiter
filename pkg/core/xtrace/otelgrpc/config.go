package otelgrpc

import (
	"context"
	"time"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Name     string
	Endpoint string
	Sampler  float64
}

// RawConfig ...
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, config); err != nil {
		xlog.Jupiter().Panic("unmarshal key", xlog.Any("err", err))
	}
	return config
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:     pkg.Name(),
		Endpoint: "localhost:4317",
		Sampler:  0,
	}
}

// Build ...
func (config *Config) Build() trace.TracerProvider {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, config.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		xlog.Jupiter().Panic("new otelgrpc", xlog.FieldMod("build"), xlog.FieldErr(err))
		return nil
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		xlog.Jupiter().Panic("new otelgrpc", xlog.FieldMod("build"), xlog.FieldErr(err))
		return nil
	}

	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.Sampler))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(traceExporter),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.TelemetrySDKLanguageGo,
			semconv.ServiceNameKey.String(config.Name),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp
}
