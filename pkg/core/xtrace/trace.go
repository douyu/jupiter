// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xtrace

import (
	"context"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
	otelOpentracing "go.opentelemetry.io/otel/bridge/opentracing"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// SetGlobalTracer ...
func SetGlobalTracer(tp trace.TracerProvider) {
	xlog.Jupiter().Info("set global tracer", xlog.FieldMod("trace"))

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, Jaeger{})

	// be compatible with opentracing
	bridge, wrapperTracerProvider := otelOpentracing.NewTracerPair(tp.Tracer(""))
	bridge.SetTextMapPropagator(propagator)
	opentracing.SetGlobalTracer(bridge)

	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(wrapperTracerProvider)
}

// Tracer is otel span tracer
type Tracer struct {
	tracer trace.Tracer
	kind   trace.SpanKind
}

// NewTracer create tracer instance
func NewTracer(kind trace.SpanKind) *Tracer {
	return &Tracer{tracer: otel.Tracer("jupiter"), kind: kind}
}

// Start start tracing span
func (t *Tracer) Start(ctx context.Context, operation string, carrier propagation.TextMapCarrier, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if (t.kind == trace.SpanKindServer || t.kind == trace.SpanKindConsumer) && carrier != nil {
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	}
	opts = append(opts, trace.WithSpanKind(t.kind))

	ctx, span := t.tracer.Start(ctx, operation, opts...)

	if (t.kind == trace.SpanKindClient || t.kind == trace.SpanKindProducer) && carrier != nil {
		otel.GetTextMapPropagator().Inject(ctx, carrier)
	}
	return ctx, span
}
