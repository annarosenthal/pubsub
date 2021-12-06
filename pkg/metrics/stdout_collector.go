package metrics

import "fmt"

type StdOutCollector struct {
}

func (s *StdOutCollector) Record(name string, tags Tags, value float64) {
	fmt.Println("record", name, tags, value)
}

func (s *StdOutCollector) Increment(name string, tags Tags) {
	fmt.Println("increment", name, tags)
}