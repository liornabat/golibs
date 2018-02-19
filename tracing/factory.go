package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

type TracingOptions struct {
	LocalHostPort  string
	ReportHostPort string
	Debug          bool
	SampleAllSpans bool
	Tracer         string
}

type Factory struct {
	name string
	opentracing.Tracer
	zipkin.Collector
	defaultTags    *Tags
	SampleAllSpans bool
	cache          *spanCache
	isClose        bool
}

type SpanOptions struct {
	Tags       *Tags
	MustSample bool
}

var tracerFactory *Factory

func InitTracing(serviceName string, ops TracingOptions, tags ...*Tags) error {

	t, c, err := InitZipkin(ops, serviceName)
	if err != nil {
		return err
	}

	tracerFactory = &Factory{
		name:           serviceName,
		Tracer:         t,
		Collector:      c,
		SampleAllSpans: ops.SampleAllSpans,
		cache:          newCache(),
	}
	if len(tags) > 0 {
		tracerFactory.defaultTags = tags[0]
	}

	return nil
}

func StartSpan(ctx context.Context, spanName string) (context.Context, *Span) {
	if ctx == nil {
		ctx = context.Background()
	}
	s, c := opentracing.StartSpanFromContext(ctx, spanName)
	span := &Span{
		s,
		c,
	}

	if span != nil {
		if tracerFactory.SampleAllSpans {
			span.SetSamplingPriority(1)
		}
	}
	ctxOut := opentracing.ContextWithSpan(ctx, span.Span)
	return ctxOut, span
}

func StartSpanFromCache(spanName string, key string, moreKeys ...string) (context.Context, *Span) {
	var span *Span
	ctxOut := context.Background()
	cacheSpan, ok := tracerFactory.cache.getSpan(key)
	if ok {
		ctx := opentracing.ContextWithSpan(context.Background(), cacheSpan.Span)
		s, c := opentracing.StartSpanFromContext(ctx, spanName)
		span = &Span{
			s,
			c,
		}

	} else {
		for _, altKey := range moreKeys {
			cacheSpan, ok := tracerFactory.cache.getSpan(altKey)
			if ok {
				ctx := opentracing.ContextWithSpan(context.Background(), cacheSpan.Span)
				s, c := opentracing.StartSpanFromContext(ctx, spanName)
				span = &Span{
					s,
					c,
				}

				break
			}
		}

	}
	if span != nil {
		if tracerFactory.SampleAllSpans {
			span.SetSamplingPriority(1)
		}
		ctxOut = opentracing.ContextWithSpan(ctxOut, span.Span)
	}

	return ctxOut, span
}

func StartSpanFromBinary(spanName string, in []byte) (context.Context, *Span, error) {
	var span *Span
	var ctxOut context.Context = context.Background()
	wireContext, err := tracerFactory.Extract(opentracing.Binary, in)
	if wireContext == nil || err != nil {
		return context.Background(), nil, err
	}
	span = &Span{
		Span: opentracing.StartSpan(spanName, ext.RPCServerOption(wireContext)),
	}

	if span != nil {
		if tracerFactory.SampleAllSpans {
			span.SetSamplingPriority(1)
		}
		ctxOut = opentracing.ContextWithSpan(context.Background(), span.Span)
	}

	return ctxOut, span, err
}

func CloseTracing() {
	tracerFactory.isClose = true
	tracerFactory.Close()

}

func GetTracer() *Factory {
	return tracerFactory
}
