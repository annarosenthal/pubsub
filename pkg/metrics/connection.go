package metrics

import (
	"errors"
	"net"
	"os"
	"time"
)

type connection struct {
	wrapped net.Conn
	listener *listener
}

func WrapConnection(c net.Conn, l *listener) net.Conn {
	return &connection{c, l}
}

func (c connection) Read(b []byte) (n int, err error) {
	n, err = c.wrapped.Read(b)
	if err !=nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		c.Close()
	}
	return n, err
}

func (c connection) Write(b []byte) (n int, err error) {
	n, err = c.wrapped.Write(b)
	if err !=nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		c.Close()
	}
	return n, err
}

func (c connection) Close() error {
	err := c.wrapped.Close()
	if err == nil {
		c.listener.connections--
		c.listener.collector.Record(ConnectionTotalMetric, EmptyTags, float64(c.listener.connections))
	}
	return err
}

func (c connection) LocalAddr() net.Addr {
	return c.wrapped.LocalAddr()
}

func (c connection) RemoteAddr() net.Addr {
	return c.wrapped.RemoteAddr()
}

func (c connection) SetDeadline(t time.Time) error {
	err := c.wrapped.SetDeadline(t)
	if err != nil {
		c.Close()
	}
	return err
}

func (c connection) SetReadDeadline(t time.Time) error {
	err := c.wrapped.SetReadDeadline(t)
	if err != nil {
		c.Close()
	}
	return err
}

func (c connection) SetWriteDeadline(t time.Time) error {
	err := c.wrapped.SetWriteDeadline(t)
	if err != nil {
		c.Close()
	}
	return err
}

