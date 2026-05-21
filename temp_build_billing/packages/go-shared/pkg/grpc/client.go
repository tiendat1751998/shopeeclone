package grpc

import (
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientOptions struct {
	Address        string
	ServiceName    string
	Timeout        time.Duration
	MaxRetries     int
	EnableCircuit  bool
}

func NewClient(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"` + roundrobin.Name + `"}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithUnaryInterceptor(otelUnaryClientInterceptor()),
		grpc.WithChainUnaryInterceptor(
			retryUnaryClientInterceptor(3),
		),
	}

	defaultOpts = append(defaultOpts, opts...)

	conn, err := grpc.Dial(addr, defaultOpts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func NewCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "grpc-circuit-breaker",
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
}
