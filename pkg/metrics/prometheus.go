package metrics

type PrometheusCollector struct {

}

func (p *PrometheusCollector) Record(name string, tags Tags, value float64) {
	panic("implement me")
}

func (p *PrometheusCollector) Increment(name string, tags Tags) {
	panic("implement me")
}
