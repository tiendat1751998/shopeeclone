package routes

import "time"

type Route struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	Methods      []string  `json:"methods"`
	ServiceName  string    `json:"service_name"`
	TargetURL    string    `json:"target_url"`
	TimeoutMs    int       `json:"timeout_ms"`
	RateLimit    int       `json:"rate_limit"`
	AuthRequired bool      `json:"auth_required"`
	Middleware   []string  `json:"middleware"`
	IsActive     bool      `json:"is_active"`
	Version      string    `json:"version"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RouteGroup struct {
	Prefix     string   `json:"prefix"`
	Routes     []Route  `json:"routes"`
	Middleware []string `json:"middleware"`
}

type RegisterRouteRequest struct {
	Path         string   `json:"path" binding:"required"`
	Methods      []string `json:"methods" binding:"required"`
	ServiceName  string   `json:"service_name" binding:"required"`
	TargetURL    string   `json:"target_url" binding:"required"`
	TimeoutMs    int      `json:"timeout_ms"`
	RateLimit    int      `json:"rate_limit"`
	AuthRequired bool     `json:"auth_required"`
	Middleware   []string `json:"middleware"`
	IsActive     bool     `json:"is_active"`
	Version      string   `json:"version"`
}

type MatchRequest struct {
	Path    string   `json:"path" binding:"required"`
	Method  string   `json:"method" binding:"required"`
	Headers map[string]string `json:"headers"`
}
