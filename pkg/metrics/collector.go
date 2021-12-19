package metrics

type Tags map[string]string

var EmptyTags = Tags{}

type Collector interface {
	Record(name string, tags Tags, value float64)
	Increment(name string, tags Tags)
}

type DefaultCollector struct {
}


func (d *DefaultCollector) Record(_ string, _ Tags, _ float64) {

}

func (d *DefaultCollector) Increment(_ string, _ Tags) {

}

const ConnectionCountMetric = "connection_count"
const SubscribeCountMetric = "subscribe_count"
const PublishCountMetric = "publish_count"
const PublishMessageSizeMetric = "publish_message_size"
const ConnectionErrorCountMetric = "connection_error_count"
const SubscribeErrorCountMetric = "subscribe_error_count"
const PublishErrorCountMetric = "publish_error_count"
const ConnectionTotalMetric = "connection_total"