package tracing

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

func InitZipkin(ops TracingOptions, serviceName string) (opentracing.Tracer, zipkin.Collector, error) {
	zipkinHTTPEndpoint := fmt.Sprintf("http://%s/api/v1/spans", ops.ReportHostPort)
	collector, err := zipkin.NewHTTPCollector(zipkinHTTPEndpoint)
	if err != nil {
		fmt.Printf("unable to create Zipkin HTTP collector: %+v\n", err)
		return nil, nil, err
	}

	// create recorder.
	recorder := zipkin.NewRecorder(collector, ops.Debug, ops.LocalHostPort, serviceName)

	// create tracer.
	t, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	if err != nil {
		fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
		return nil, nil, err
	}

	// explicitly set our tracer to be the default tracer.
	opentracing.SetGlobalTracer(t)

	return t, collector, nil

}
