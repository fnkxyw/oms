package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	orderLabel   = "operation"
	handlerLabel = "handler"
	codeLabel    = "code"
)

var (
	orderTotalOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "total_number_of_ok_orders_operations",
		Help: "total number of notified positions by contact total",
	}, []string{
		orderLabel,
	})

	okRespByHandlerTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pup_ok_response_by_handler_total",
		Help: "total number of ok responses in handler total",
	}, []string{
		handlerLabel,
	})

	badRespByHandlerTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pup_bad_response_by_handler_total",
		Help: "total number of bad responses in handler total with resp code",
	}, []string{
		handlerLabel,
	})
)

func IncOrderTotalOperations(operation string) {
	orderTotalOperations.With(prometheus.Labels{
		orderLabel: operation,
	}).Add(float64(1))
}

func IncOkRespByHandler(handler string) {
	okRespByHandlerTotal.With(prometheus.Labels{
		handlerLabel: handler,
	}).Inc()
}

func IncBadRespByHandler(handler string) {
	badRespByHandlerTotal.With(prometheus.Labels{
		handlerLabel: handler,
	}).Inc()
}
