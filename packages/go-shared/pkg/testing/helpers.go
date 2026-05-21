package testing

import (
	"context"
	"net"
	"time"
)

type TestLogger struct {
	Logf func(format string, args ...interface{})
}

func WaitForPort(ctx context.Context, addr string, timeout time.Duration) error {
	deadline := time.After(timeout)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline:
			return nil
		default:
			conn, err := net.DialTimeout("tcp", addr, time.Second)
			if err == nil {
				conn.Close()
				return nil
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type NullWriter struct{}

func (n NullWriter) Write(p []byte) (int, error) { return len(p), nil }
