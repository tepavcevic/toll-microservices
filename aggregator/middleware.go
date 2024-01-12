package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/tepavcevic/toll-microservices/types"
)

type MetricsMiddleware struct {
	next           Aggregator
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
}

type LogMiddleware struct {
	next Aggregator
}

func NewMetricsMiddleware(n Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calculator_request_counter",
		Name:      "calculate",
	})
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "calculator_error_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "calculator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           n,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(started time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(started).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(dist)
	return err
}

func (m *MetricsMiddleware) GetInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(started time.Time) {
		m.reqLatencyCalc.Observe(float64(time.Since(started).Seconds()))
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())

	invoice, err = m.next.GetInvoice(obuID)
	return
}

func NewLogMiddleware(n Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: n,
	}
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(started time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"took":     time.Since(started),
			"obuID":    dist.OBUID,
			"distance": dist.Value,
		}).Info("aggregating distance")
	}(time.Now())

	err = m.next.AggregateDistance(dist)
	return
}

func (m *LogMiddleware) GetInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(started time.Time) {
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"obuID":         obuID,
			"totalDistance": invoice.TotalDistance,
			"totalAmount":   invoice.TotalAmount,
			"took":          time.Since(started),
		}).Info("getting the aggregated distance by obuID")
	}(time.Now())

	invoice, err = m.next.GetInvoice(obuID)
	return
}
