package tracer

import (
	"context"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/logger"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

func MustSetup(ctx context.Context, serviceName string) {
	cfg := traceconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &traceconfig.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(traceconfig.Logger(jaeger.StdLogger), traceconfig.Metrics(prometheus.New()))
	if err != nil {
		logger.Errorf(ctx, "ERROR: cannot init Jaeger %s", err)
		return
	}

	go func() {
		onceCloser := sync.OnceFunc(func() {
			logger.Infof(ctx, "closing tracer")
			if err = closer.Close(); err != nil {
				logger.Errorf(ctx, "error closing tracer: %s", err)
				return
			}
		})

		for {
			<-ctx.Done()
			onceCloser()
		}
	}()

	opentracing.SetGlobalTracer(tracer)
}
