package traffic



type MatchCondition struct {
	Headers map[string]string `json:"headers"`
	PathPrefix string         `json:"path_prefix"`
	Methods   []string        `json:"methods"`
}

type TrafficRule struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	SourceService    string          `json:"source_service"`
	DestinationService string        `json:"destination_service"`
	MatchConditions  MatchCondition  `json:"match_conditions"`
	Weight           int             `json:"weight"`
	MirrorPercentage int             `json:"mirror_percentage"`
}

type VirtualService struct {
	Name  string        `json:"name"`
	Rules []*TrafficRule `json:"rules"`
}

type CircuitBreakerSettings struct {
	MaxConnections     int `json:"max_connections"`
	MaxPendingRequests int `json:"max_pending_requests"`
	MaxRequests        int `json:"max_requests"`
	MaxRetries         int `json:"max_retries"`
}

type ConnectionPoolSettings struct {
	TCP          TCPPoolSettings  `json:"tcp"`
	HTTP         HTTPPoolSettings `json:"http"`
}

type TCPPoolSettings struct {
	MaxConnections int `json:"max_connections"`
	ConnectTimeout int `json:"connect_timeout_ms"`
}

type HTTPPoolSettings struct {
	HTTP1MaxPendingRequests  int `json:"http1_max_pending_requests"`
	HTTP2MaxRequests         int `json:"http2_max_requests"`
	MaxRequestsPerConnection int `json:"max_requests_per_connection"`
}

type OutlierDetectionSettings struct {
	ConsecutiveErrors    int `json:"consecutive_errors"`
	Interval             int `json:"interval_ms"`
	BaseEjectionTime     int `json:"base_ejection_time_ms"`
	MaxEjectionPercent   int `json:"max_ejection_percent"`
}

type DestinationRule struct {
	Name                   string                   `json:"name"`
	TrafficPolicy          CircuitBreakerSettings   `json:"traffic_policy"`
	ConnectionPool         ConnectionPoolSettings   `json:"connection_pool"`
	OutlierDetection       OutlierDetectionSettings `json:"outlier_detection"`
}
