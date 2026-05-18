package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func hasInfrastructure() bool {
	return os.Getenv("MYSQL_HOST") != "" && os.Getenv("REDIS_ADDR") != ""
}

func setupE2EEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	return engine
}

func TestE2E_RegisterUser(t *testing.T) {
	if !hasInfrastructure() {
		t.Skip("Skipping E2E test: MYSQL_HOST and REDIS_ADDR not set")
	}

	engine := setupE2EEngine()

	// Register
	registerBody := map[string]string{
		"email":            "e2e-test@example.com",
		"password":         "StrongP@ss1",
		"confirm_password": "StrongP@ss1",
		"username":         "e2euser",
		"display_name":     "E2E User",
	}
	body, _ := json.Marshal(registerBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var tokens map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&tokens); err != nil {
		t.Fatal(err)
	}
	accessToken, _ := tokens["access_token"].(string)
	refreshToken, _ := tokens["refresh_token"].(string)
	if accessToken == "" || refreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}

	// Get profile
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("profile: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var profile map[string]interface{}
	json.NewDecoder(w.Body).Decode(&profile)
	if profile["email"] != "e2e-test@example.com" {
		t.Errorf("expected email e2e-test@example.com, got %v", profile["email"])
	}

	// Refresh token
	refreshBody := map[string]string{"refresh_token": refreshToken}
	body, _ = json.Marshal(refreshBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var newTokens map[string]interface{}
	json.NewDecoder(w.Body).Decode(&newTokens)
	newAccessToken, _ := newTokens["access_token"].(string)
	newRefreshToken, _ := newTokens["refresh_token"].(string)
	if newAccessToken == "" || newRefreshToken == "" {
		t.Fatal("expected new access and refresh tokens after refresh")
	}
	if newAccessToken == accessToken {
		t.Error("access token should rotate after refresh")
	}

	// Old refresh token should be invalidated (rotation)
	body, _ = json.Marshal(refreshBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("old refresh token should be rejected after rotation")
	}

	// Reuse detection: try reusing the OLD refresh token again
	body, _ = json.Marshal(refreshBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("refresh token reuse should trigger revocation")
	}

	// Logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+newAccessToken)
	req.Header.Set("X-Refresh-Token", newRefreshToken)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("logout: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Access token should be blacklisted after logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/profile", nil)
	req.Header.Set("Authorization", "Bearer "+newAccessToken)
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("profile after logout: expected 401, got %d", w.Code)
	}
}

func TestE2E_LoginFlow(t *testing.T) {
	if !hasInfrastructure() {
		t.Skip("Skipping E2E test: MYSQL_HOST and REDIS_ADDR not set")
	}

	engine := setupE2EEngine()

	// Register
	registerBody := map[string]string{
		"email":            "e2e-login@example.com",
		"password":         "StrongP@ss1",
		"confirm_password": "StrongP@ss1",
		"username":         "e2elogin",
	}
	body, _ := json.Marshal(registerBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d", w.Code)
	}

	// Login with valid credentials
	loginBody := map[string]string{
		"email":    "e2e-login@example.com",
		"password": "StrongP@ss1",
	}
	body, _ = json.Marshal(loginBody)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var tokens map[string]interface{}
	json.NewDecoder(w.Body).Decode(&tokens)
	if tokens["access_token"] == "" || tokens["refresh_token"] == "" {
		t.Fatal("expected access and refresh tokens")
	}
}

func TestE2E_LoginValidation(t *testing.T) {
	if !hasInfrastructure() {
		t.Skip("Skipping E2E test: MYSQL_HOST and REDIS_ADDR not set")
	}

	engine := setupE2EEngine()

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{"missing fields", map[string]string{}, http.StatusBadRequest},
		{"wrong password", map[string]string{"email": "nonexistent@example.com", "password": "wrong"}, http.StatusUnauthorized},
		{"empty email", map[string]string{"email": "", "password": "somepass"}, http.StatusBadRequest},
		{"empty password", map[string]string{"email": "a@b.com", "password": ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestE2E_RegistrationValidation(t *testing.T) {
	if !hasInfrastructure() {
		t.Skip("Skipping E2E test: MYSQL_HOST and REDIS_ADDR not set")
	}

	engine := setupE2EEngine()

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{"missing fields", map[string]string{}, http.StatusBadRequest},
		{"weak password", map[string]string{"email": "a@b.com", "password": "short", "username": "test"}, http.StatusUnprocessableEntity},
		{"mismatched passwords", map[string]string{"email": "a@b.com", "password": "StrongP@ss1", "confirm_password": "DifferentP@ss1", "username": "test"}, http.StatusBadRequest},
		{"missing username", map[string]string{"email": "a@b.com", "password": "StrongP@ss1"}, http.StatusBadRequest},
		{"invalid email", map[string]string{"email": "invalid", "password": "StrongP@ss1", "username": "test"}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("expected %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestE2E_UnauthenticatedAccess(t *testing.T) {
	if !hasInfrastructure() {
		t.Skip("Skipping E2E test: MYSQL_HOST and REDIS_ADDR not set")
	}

	engine := setupE2EEngine()

	protectedEndpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/auth/profile"},
		{http.MethodGet, "/api/v1/auth/sessions"},
		{http.MethodPost, "/api/v1/auth/logout"},
		{http.MethodPost, "/api/v1/auth/logout-all"},
	}

	for _, ep := range protectedEndpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(ep.method, ep.path, nil)
			req.Header.Set("Content-Type", "application/json")
			engine.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", w.Code)
			}
		})
	}
}
