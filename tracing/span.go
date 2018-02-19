package tracing

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

// Span Struct
type Span struct {
	opentracing.Span

	context.Context
}

func (s *Span) StoreSpanToCache(key string) *Span {
	tracerFactory.cache.putSpan(key, s)
	return s
}

func (s *Span) SetBaggage(key, value string) *Span {
	s.SetBaggageItem(key, value)
	return s
}

func (s *Span) GetBaggage(key string) string {
	return s.Span.BaggageItem(key)
}

func (s *Span) SetError(err error) *Span {
	s.Span.LogKV("error-message", err.Error())
	s.Span.LogKV("event", "error")
	s.SetErrorTag()
	return s
}

func (s *Span) Log(message string) *Span {
	s.Span.LogKV("message", message)
	return s
}

func (s *Span) LogKV(key, value string) *Span {
	s.Span.LogKV(key, value)
	return s
}

func (s *Span) Finish() {
	if tracerFactory.isClose {
		panic("tracer is closed")
	}
	s.Span.Finish()
}

// General Tags
func (s *Span) SetComponent(value string) *Span {
	s.Span.SetTag("component", value)
	return s
}

func (s *Span) MessageBusDestination(value string) *Span {
	s.Span.SetTag("message_bus.destination", value)
	return s
}

func (s *Span) SetErrorTag() *Span {
	s.Span.SetTag("error", true)
	return s
}

// Peer Tags
func (s *Span) SetSamplingPriority(value uint16) *Span {
	s.Span.SetTag("sampling.priority", value)
	return s
}

func (s *Span) SetPeerService(value string) *Span {
	s.Span.SetTag("peer.service", value)
	return s
}

func (s *Span) SetPeerAddress(value string) *Span {
	s.Span.SetTag("peer.address", value)
	return s
}

func (s *Span) SetPeerHostname(value string) *Span {
	s.Span.SetTag("peer.hostname", value)
	return s
}

func (s *Span) SetPeerHostIPv4(value string) *Span {
	s.Span.SetTag("peer.ipv4", value)
	return s
}

func (s *Span) SetPeerHostIPv6(value string) *Span {
	s.Span.SetTag("peer.ipv6", value)
	return s
}

func (s *Span) SetPeerPort(value uint16) *Span {
	s.Span.SetTag("peer.port", value)
	return s
}

// HTTP Tags

func (s *Span) SetHTTPUrl(value string) *Span {
	s.Span.SetTag("http.url", value)
	return s
}

func (s *Span) SetHTTPMethod(value string) *Span {
	s.Span.SetTag("http.method", value)
	return s
}

func (s *Span) SetHTTPStatusCode(value uint16) *Span {
	s.Span.SetTag("http.status_code", value)
	return s
}

// DB Tags

func (s *Span) SetDBInstance(value string) *Span {
	s.Span.SetTag("db.instance", value)
	return s
}

func (s *Span) SetDBStatement(value string) *Span {
	s.Span.SetTag("db.statement", value)
	return s
}

func (s *Span) SetDBType(value string) *Span {
	s.Span.SetTag("db.type", value)
	return s
}

func (s *Span) SetDBUser(value string) *Span {
	s.Span.SetTag("db.user", value)
	return s
}

// Set SpanKind

func (s *Span) SetSpanKindRPCClient() *Span {
	s.Span.SetTag("span.kind", "client")
	return s
}

func (s *Span) SetSpanKindRPCServer() *Span {
	s.Span.SetTag("span.kind", "server")
	return s
}

func (s *Span) SetSpanKindProducer() *Span {
	s.Span.SetTag("span.kind", "producer")
	return s
}

func (s *Span) SetSpanKindConsumer() *Span {
	s.Span.SetTag("span.kind", "consumer")
	return s
}

func (s *Span) SetTag(key string, value interface{}) *Span {
	s.Span.SetTag(key, value)
	return s
}

func (s *Span) ToBinary(out []byte) {
	tracerFactory.Inject(s.Span.Context(), opentracing.Binary, out)
	return
}
