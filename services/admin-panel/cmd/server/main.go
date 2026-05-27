package main

import (
	"embed"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	automaxprocs "go.uber.org/automaxprocs/maxprocs"
)

//go:embed web/**
var webFS embed.FS

var logger *zap.Logger

func init() {
	// Tune GC for low-latency: more frequent GCs, less heap growth
	if gogc := os.Getenv("GOGC"); gogc == "" {
		os.Setenv("GOGC", "50")
	}
}

func main() {
	log, _ := zap.NewProduction()
	// Auto-tune GOMAXPROCS for container environments
	_, _ = automaxprocs.Set()

	logger = log

	port := getEnv("ADMIN_PORT", "3001")
	apiGateway := getEnv("API_GATEWAY_URL", "http://localhost:8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/", authProxyHandler(apiGateway))
	mux.HandleFunc("/api/admin/", adminProxyHandler(apiGateway))
	mux.HandleFunc("/api/", apiProxyHandler(apiGateway))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"admin-panel"}`))
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ready"}`))
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("# Admin Panel metrics\n"))
	})
	mux.HandleFunc("/", spaHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMiddleware(loggingMiddleware(mux)),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("Admin Panel starting", zap.String("port", port), zap.String("gateway", apiGateway))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Admin server failed", zap.Error(err))
	}
}

func spaHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/health" || r.URL.Path == "/ready" || r.URL.Path == "/metrics" {
		return
	}

	if strings.Contains(r.URL.Path, ".") && r.URL.Path != "/" {
		data, err := webFS.ReadFile("web/" + strings.TrimPrefix(r.URL.Path, "/"))
		if err == nil {
			ext := r.URL.Path[strings.LastIndex(r.URL.Path, ".")+1:]
			ct := "text/plain"
			switch ext {
			case "css":
				ct = "text/css"
			case "js":
				ct = "application/javascript"
			case "png":
				ct = "image/png"
			case "svg":
				ct = "image/svg+xml"
			case "ico":
				ct = "image/x-icon"
			}
			w.Header().Set("Content-Type", ct)
			w.Write(data)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := webFS.ReadFile("web/index.html")
	if err != nil {
		logger.Error("Failed to load index.html", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(tmpl)
}

func apiProxyHandler(gateway string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized","message":"Authorization header required"}`))
			return
		}
		proxyTo(w, r, gateway+r.URL.Path, authHeader)
	}
}

func adminProxyHandler(gateway string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		targetURL := gateway + "/api/v1/admin" + strings.TrimPrefix(r.URL.Path, "/api/admin")
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, `{"error":"bad gateway"}`, http.StatusBadGateway)
			return
		}
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Admin-Request", "true")
		doProxy(w, req)
	}
}

func authProxyHandler(gateway string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL := gateway + r.URL.Path
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}
		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, `{"error":"bad gateway"}`, http.StatusBadGateway)
			return
		}
		for k, v := range r.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
		doProxy(w, req)
	}
}

func proxyTo(w http.ResponseWriter, r *http.Request, targetURL, authHeader string) {
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, `{"error":"bad gateway"}`, http.StatusBadGateway)
		return
	}
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	doProxy(w, req)
}

func doProxy(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		if strings.ToLower(k) != "transfer-encoding" {
			for _, vv := range v {
				w.Header().Add(k, vv)
			}
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Requested-With")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Debug("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
