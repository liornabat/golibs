package tracing

import (
	"os"
	"testing"
	"time"
	//	"errors"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

func func2(ctx context.Context, cnt int) {

	span := NewSpanFromContext(fmt.Sprintf("func2-%d", cnt), ctx)
	defer span.Finish()
	span.SetSamplingPriority(1)
	time.Sleep(1 * time.Second)
	span.SetDBType("dasdsa")
	span.Log("done")
	span.LogKV("baggage", span.GetBaggage("userId"))
	//span.SetError(errors.New("error tag"))
}

func TestNewSpanFactory(t *testing.T) {
	require := require.New(t)
	err := InitTracing("root", TracingOptions{
		LocalHostPort:  "127.0.0.1:0",
		Debug:          true,
		ReportHostPort: "localhost:9412",
	})
	require.NoError(err)
	defer CloseTracing()

	span := NewRootSpan("span1")
	span.SetSpanKindRPCServer()
	span.SetSamplingPriority(1)
	span.LogKV("asd", "asdasd")
	span.Log("message1")
	span.SetTag("userId", 1111)
	span.SetBaggage("userId", "1111")
	for i := 0; i < 10; i++ {
		go func2(span.GetContextFromSpan(), i)
	}

	time.Sleep(10 * time.Second)
	span.LogKV("baggage", span.GetBaggage("userId"))
	span.Log("done")
	span.Finish()
}
