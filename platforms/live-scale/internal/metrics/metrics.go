package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	SFUNodesTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livescale_sfu_nodes_total",
		Help: "Total number of registered SFU nodes",
	})

	SFUStreamsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livescale_sfu_streams_total",
		Help: "Total number of active streams across all SFU nodes",
	})

	StreamHealthStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livescale_stream_health_status",
		Help: "Stream health status: 1=healthy, 0=degraded, -1=down",
	}, []string{"stream_id"})

	CDNPurgeTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "livescale_cdn_purge_total",
		Help: "Total number of CDN purge requests",
	})

	CDNEndpointLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "livescale_cdn_endpoint_latency_ms",
		Help:    "CDN endpoint latency in milliseconds",
		Buckets: []float64{5, 10, 25, 50, 100, 200, 500},
	})

	RegionLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "livescale_region_latency_ms",
		Help:    "Inter-region latency in milliseconds",
		Buckets: []float64{10, 25, 50, 100, 200, 500, 1000},
	}, []string{"from_region", "to_region"})

	TranscodeJobsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "livescale_transcode_jobs_total",
		Help: "Total number of transcode jobs created",
	})

	TranscodeJobsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livescale_transcode_jobs_active",
		Help: "Number of active transcode jobs",
	})

	WSClusterNodesTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livescale_ws_cluster_nodes_total",
		Help: "Total number of WebSocket cluster nodes",
	})

	WSRoomsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "livescale_ws_rooms_total",
		Help: "Total number of rooms across WebSocket cluster",
	})
)
