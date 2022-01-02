package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
)

type PrometheusCollector struct {
	server   *http.Server
	counters map[string]*prometheus.CounterVec
	gauges   map[string]*prometheus.GaugeVec
}

func NewPrometheusCollector() *PrometheusCollector {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Handler: mux}
	counters := map[string]*prometheus.CounterVec{
		ConnectionCountMetric:      newCounter(ConnectionCountMetric, "useful"),
		SubscribeCountMetric:       newCounter(SubscribeCountMetric, "useful", "topic"),
		PublishCountMetric:         newCounter(PublishCountMetric, "useful", "topic"),
		ConnectionErrorCountMetric: newCounter(ConnectionErrorCountMetric, "useful"),
		SubscribeErrorCountMetric:  newCounter(SubscribeErrorCountMetric, "useful", "topic"),
		PublishErrorCountMetric:    newCounter(PublishErrorCountMetric, "useful", "topic"),
	}
	gauges := map[string]*prometheus.GaugeVec{
		PublishMessageSizeMetric: newGauge(PublishMessageSizeMetric, "useful", "topic"),
		ConnectionTotalMetric:    newGauge(ConnectionTotalMetric, "useful"),
	}
	return &PrometheusCollector{server, counters, gauges}
}

func (p *PrometheusCollector) Start() error {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		return err
	}
	go p.server.Serve(listener)
	return nil
}

func (p *PrometheusCollector) Close() error {
	return p.server.Close()
}

func (p *PrometheusCollector) Record(name string, tags Tags, value float64) {
	gauge, ok := p.gauges[name]
	if !ok {
		return
	}
	gauge.With(prometheus.Labels(tags)).Set(value)
}

func (p *PrometheusCollector) Increment(name string, tags Tags) {
	counter, ok := p.counters[name]
	if !ok {
		return
	}
	counter.With(prometheus.Labels(tags)).Inc()
}

func newCounter(name string, help string, labelNames ...string) *prometheus.CounterVec {
	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: help,
		},
		labelNames)
}

func newGauge(name string, help string, labelNames ...string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		},
		labelNames)
}
