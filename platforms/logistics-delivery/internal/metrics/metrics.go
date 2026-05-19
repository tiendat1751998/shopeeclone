package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ShipmentCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_shipment_created_total",
		Help: "Total number of shipments created",
	})

	ShipmentDeliveredTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_shipment_delivered_total",
		Help: "Total number of shipments delivered",
	})

	ShipmentFailedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_shipment_failed_total",
		Help: "Total number of shipments failed",
	})

	ShipmentProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "logistics_shipment_processing_duration_seconds",
		Help:    "Shipment processing duration in seconds",
		Buckets: prometheus.DefBuckets,
	})

	TrackingEventTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_tracking_event_total",
		Help: "Total number of tracking events",
	})

	TrackingLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "logistics_tracking_event_latency_seconds",
		Help:    "Tracking event latency in seconds",
		Buckets: prometheus.DefBuckets,
	})

	DispatchDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "logistics_dispatch_duration_seconds",
		Help:    "Dispatch operation duration in seconds",
		Buckets: prometheus.DefBuckets,
	})

	CourierAPILatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "logistics_courier_api_latency_seconds",
		Help:    "Courier API call latency in seconds",
		Buckets: prometheus.DefBuckets,
	})

	ETACalculationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "logistics_eta_calculation_duration_seconds",
		Help:    "ETA calculation duration in seconds",
		Buckets: prometheus.DefBuckets,
	})

	DeliveryDelayCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_delivery_delay_total",
		Help: "Total number of delayed deliveries",
	})

	ReplayProcessedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistics_replay_processed_total",
		Help: "Total number of replay events processed",
	})

	ActiveCouriers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logistics_active_couriers",
		Help: "Number of active couriers",
	})

	ActiveShipments = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "logistics_active_shipments",
		Help: "Number of active shipments",
	})
)

func RecordShipmentProcessingDuration(d time.Duration) {
	ShipmentProcessingDuration.Observe(d.Seconds())
}

func RecordTrackingEventLatency(d time.Duration) {
	TrackingLatency.Observe(d.Seconds())
}

func RecordDispatchDuration(d time.Duration) {
	DispatchDuration.Observe(d.Seconds())
}

func RecordCourierAPILatency(d time.Duration) {
	CourierAPILatency.Observe(d.Seconds())
}

func RecordETACalculationDuration(d time.Duration) {
	ETACalculationDuration.Observe(d.Seconds())
}
