package tracing

import (
	"github.com/opentracing/opentracing-go"
)

type Tags struct {
	m map[string]interface{}
}

func NewTags() *Tags {
	t := &Tags{
		m: make(map[string]interface{}),
	}
	return t
}

// General Tags
func (t *Tags) SetComponent(value string) *Tags {
	t.m["component"] = value
	return t
}

func (t *Tags) MessageBusDestination(value string) *Tags {
	t.m["message_bus.destination"] = value
	return t
}

func (t *Tags) SetError() *Tags {
	t.m["error"] = true
	return t
}

// Peer Tags
func (t *Tags) SetSamplingPriority(value uint16) *Tags {
	t.m["sampling.priority"] = value
	return t
}

func (t *Tags) SetPeerService(value string) *Tags {
	t.m["peer.service"] = value
	return t
}

func (t *Tags) SetPeerAddress(value string) *Tags {
	t.m["peer.address"] = value
	return t
}

func (t *Tags) SetPeerHostname(value string) *Tags {
	t.m["peer.hostname"] = value
	return t
}

func (t *Tags) SetPeerHostIPv4(value string) *Tags {
	t.m["peer.ipv4"] = value
	return t
}

func (t *Tags) SetPeerHostIPv6(value string) *Tags {
	t.m["peer.ipv6"] = value
	return t
}

func (t *Tags) SetPeerPort(value uint16) *Tags {
	t.m["peer.port"] = value
	return t
}

// HTTP Tags

func (t *Tags) SetHTTPUrl(value string) *Tags {
	t.m["http.url"] = value
	return t
}

func (t *Tags) SetHTTPMethod(value string) *Tags {
	t.m["http.method"] = value
	return t
}

func (t *Tags) SetHTTPStatusCode(value uint16) *Tags {
	t.m["http.status_code"] = value
	return t
}

// DB Tags

func (t *Tags) SetDBInstance(value string) *Tags {
	t.m["db.instance"] = value
	return t
}

func (t *Tags) SetDBStatement(value string) *Tags {
	t.m["db.statement"] = value
	return t
}

func (t *Tags) SetDBType(value string) *Tags {
	t.m["db.type"] = value
	return t
}

func (t *Tags) SetDBUser(value string) *Tags {
	t.m["db.user"] = value
	return t
}

// Set SpanKind

func (t *Tags) SetSpanKindRPCClient() *Tags {
	t.m["span.kind"] = "client"
	return t
}

func (t *Tags) SetSpanKindRPCServer() *Tags {
	t.m["span.kind"] = "server"
	return t
}

func (t *Tags) SetSpanKindProducer() *Tags {
	t.m["span.kind"] = "producer"
	return t
}

func (t *Tags) SetSpanKindConsumer() *Tags {
	t.m["span.kind"] = "consumer"
	return t
}

func (t *Tags) SetTag(key string, value interface{}) *Tags {
	t.m[key] = value
	return t
}

func (t *Tags) SetTagsToSpan(span opentracing.Span) opentracing.Span {
	for key, value := range t.m {
		span.SetTag(key, value)
	}

	return span
}
