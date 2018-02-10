package tracing

import (
	"github.com/opentracing/opentracing-go"

	"context"

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

// func NewRootSpan(name string, ops ...SpanOptions) (span *Span) {

// 	span = &Span{
// 		Span: tracerFactory.Tracer.StartSpan(name),
// 	}
// 	if len(ops) > 0 {
// 		if (ops[0].Tags) != nil {
// 			span.Span = ops[0].Tags.SetTagsToSpan(span.Span)
// 		}
// 		if ops[0].MustSample {
// 			span.SetSamplingPriority(1)
// 		}

// 	}
// 	if tracerFactory.SampleAllSpans {
// 		span.SetSamplingPriority(1)
// 	}
// 	return
// }

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

// func NewSpanFromContext(name string, ctx context.Context, ops ...SpanOptions) (span *Span) {
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	s, c := opentracing.StartSpanFromContext(ctx, name)
// 	span = &Span{
// 		s,
// 		c,
// 	}

// 	if len(ops) > 0 {
// 		if (ops[0].Tags) != nil {
// 			span.Span = ops[0].Tags.SetTagsToSpan(span.Span)
// 		}
// 		if ops[0].MustSample {
// 			span.SetSamplingPriority(1)
// 		}

// 	}
// 	if tracerFactory.SampleAllSpans {
// 		span.SetSamplingPriority(1)
// 	}

// 	return
// }
func StartSpanFromCache(spanName string, key string, moreKeys ...string) (ctx context.Context, span *Span) {

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
	}

	return
}

// func NewSpanFromCache(name string, key string, moreKeys ...string) ( span *Span) {

// 	cacheSpan, ok := tracerFactory.cache.getSpan(key)
// 	if ok {
// 		cacheCtx := cacheSpan.GetContextFromSpan()
// 		s, c := opentracing.StartSpanFromContext(cacheCtx, name)
// 		span = &Span{
// 			s,
// 			c,
// 		}

// 	} else {
// 		for _, altKey := range moreKeys {
// 			cacheSpan, ok := tracerFactory.cache.getSpan(altKey)
// 			if ok {
// 				cacheCtx := cacheSpan.GetContextFromSpan()
// 				s, c := opentracing.StartSpanFromContext(cacheCtx, name)
// 				span = &Span{
// 					s,
// 					c,
// 				}

// 				break
// 			}
// 		}

// 	}
// 	if span != nil {
// 		if tracerFactory.SampleAllSpans {
// 			span.SetSamplingPriority(1)
// 		}
// 	}

// 	return
// }

// func StoreSpanToCache(key string, span *Span) {
// 	tracerFactory.cache.putSpan(key, span)

// }

func CloseTracing() {
	tracerFactory.Close()
}

func GetTracer() *Factory {
	return tracerFactory
}
