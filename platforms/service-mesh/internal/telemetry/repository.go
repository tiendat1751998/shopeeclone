package telemetry

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ServiceGraph struct {
	Nodes []ServiceNode `json:"nodes"`
	Edges []ServiceEdge `json:"edges"`
}

type ServiceNode struct {
	Name      string `json:"name"`
	CallCount int64  `json:"call_count"`
}

type ServiceEdge struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	CallCount   int64   `json:"call_count"`
	ErrorRate   float64 `json:"error_rate"`
	P50Latency  float64 `json:"p50_latency"`
	P99Latency  float64 `json:"p99_latency"`
}

type TraceSpan struct {
	TraceID      string            `json:"trace_id"`
	SpanID       string            `json:"span_id"`
	ParentSpanID string            `json:"parent_span_id"`
	Service      string            `json:"service"`
	Operation    string            `json:"operation"`
	StartTime    time.Time         `json:"start_time"`
	DurationMs   float64           `json:"duration_ms"`
	Status       string            `json:"status"`
	Tags         map[string]string `json:"tags"`
}

type Repository interface {
	RecordCall(ctx context.Context, source, destination string, durationMs float64, status string) error
	GetServiceGraph(ctx context.Context) (*ServiceGraph, error)
	RecordSpan(ctx context.Context, span *TraceSpan) error
	GetTraces(ctx context.Context, service string, limit int) ([]*TraceSpan, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	edges    map[string]*ServiceEdge
	spans    []*TraceSpan
	spanMu   sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		edges: make(map[string]*ServiceEdge),
		spans: make([]*TraceSpan, 0, 1000),
	}
}

func edgeKey(source, destination string) string {
	return source + "->" + destination
}

func (r *InMemoryRepository) RecordCall(ctx context.Context, source, destination string, durationMs float64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := edgeKey(source, destination)
	edge, ok := r.edges[key]
	if !ok {
		edge = &ServiceEdge{
			Source:      source,
			Destination: destination,
		}
		r.edges[key] = edge
	}

	edge.CallCount++

	durations := edge.CallCount
	edge.P50Latency = edge.P50Latency + (durationMs-edge.P50Latency)/float64(durations)
	edge.P99Latency = edge.P99Latency + (durationMs-edge.P99Latency)/float64(durations)

	if status != "ok" && status != "success" {
		edge.CallCount++
		edge.ErrorRate = float64(edge.CallCount) / float64(durations)
	} else {
		if edge.ErrorRate > 0 {
			edge.ErrorRate = edge.ErrorRate * (1 - 1/float64(durations))
		}
	}
	return nil
}

func (r *InMemoryRepository) GetServiceGraph(ctx context.Context) (*ServiceGraph, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	nodesMap := make(map[string]*ServiceNode)
	var edges []ServiceEdge

	for _, edge := range r.edges {
		edges = append(edges, *edge)
		if _, ok := nodesMap[edge.Source]; !ok {
			nodesMap[edge.Source] = &ServiceNode{Name: edge.Source}
		}
		if _, ok := nodesMap[edge.Destination]; !ok {
			nodesMap[edge.Destination] = &ServiceNode{Name: edge.Destination}
		}
		nodesMap[edge.Source].CallCount += edge.CallCount
		nodesMap[edge.Destination].CallCount += edge.CallCount
	}

	var nodes []ServiceNode
	for _, n := range nodesMap {
		nodes = append(nodes, *n)
	}

	return &ServiceGraph{Nodes: nodes, Edges: edges}, nil
}

func (r *InMemoryRepository) RecordSpan(ctx context.Context, span *TraceSpan) error {
	if span.TraceID == "" {
		span.TraceID = uuid.New().String()
	}
	if span.SpanID == "" {
		span.SpanID = uuid.New().String()
	}
	if span.StartTime.IsZero() {
		span.StartTime = time.Now()
	}

	r.spanMu.Lock()
	defer r.spanMu.Unlock()

	if len(r.spans) >= 10000 {
		r.spans = r.spans[len(r.spans)-5000:]
	}
	r.spans = append(r.spans, span)
	return nil
}

func (r *InMemoryRepository) GetTraces(ctx context.Context, service string, limit int) ([]*TraceSpan, error) {
	r.spanMu.RLock()
	defer r.spanMu.RUnlock()

	if limit <= 0 {
		limit = 100
	}

	var result []*TraceSpan
	for i := len(r.spans) - 1; i >= 0 && len(result) < limit; i-- {
		span := r.spans[i]
		if service == "" || span.Service == service {
			result = append(result, span)
		}
	}

	if result == nil {
		result = []*TraceSpan{}
	}
	return result, nil
}

func (r *InMemoryRepository) GetDependencies(ctx context.Context) (map[string][]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	deps := make(map[string][]string)
	for _, edge := range r.edges {
		deps[edge.Source] = append(deps[edge.Source], edge.Destination)
	}
	return deps, nil
}
