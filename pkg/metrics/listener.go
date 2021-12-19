package metrics

import "net"

type listener struct {
	wrapped net.Listener
	collector Collector
	connections uint16
}

func WrapListener(l net.Listener, collector Collector) net.Listener {
	return &listener{l, collector, 0}
}

// track if a connection is solved
func (l *listener) Accept() (net.Conn, error) {
	conn, err := l.wrapped.Accept()
	if err != nil {
		l.collector.Increment(ConnectionErrorCountMetric, EmptyTags)
	} else {
		l.collector.Increment(ConnectionCountMetric, EmptyTags)
		l.connections++
		l.collector.Record(ConnectionTotalMetric, EmptyTags, float64(l.connections))
	}
	return WrapConnection(conn, l), err
}

func (l *listener) Close() error {
	return l.wrapped.Close()
}

func (l *listener) Addr() net.Addr {
	return l.wrapped.Addr()
}
