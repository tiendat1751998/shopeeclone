package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type APIKey struct {
	Key       string    `json:"key"`
	Service   string    `json:"service"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type JWTClaims struct {
	Subject   string            `json:"sub"`
	Issuer    string            `json:"iss"`
	ExpiresAt int64             `json:"exp"`
	IssuedAt  int64             `json:"iat"`
	Extra     map[string]string `json:"extra,omitempty"`
}

type APIKeyStore struct {
	mu   sync.RWMutex
	keys map[string]*APIKey
}

func NewAPIKeyStore() *APIKeyStore {
	return &APIKeyStore{
		keys: make(map[string]*APIKey),
	}
}

func (s *APIKeyStore) Store(key *APIKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[key.Key] = key
	return nil
}

func (s *APIKeyStore) Get(key string) (*APIKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	k, ok := s.keys[key]
	if !ok {
		return nil, nil
	}
	return k, nil
}

func (s *APIKeyStore) List() ([]*APIKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*APIKey, 0, len(s.keys))
	for _, k := range s.keys {
		result = append(result, k)
	}
	return result, nil
}

type JWTHandler struct {
	secret []byte
}

func NewJWTHandler(secret string) *JWTHandler {
	return &JWTHandler{secret: []byte(secret)}
}

func (j *JWTHandler) Sign(claims JWTClaims) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claims.IssuedAt = time.Now().Unix()
	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := headerB64 + "." + claimsB64
	sig := j.sign(signingInput)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return signingInput + "." + sigB64, nil
}

func (j *JWTHandler) Verify(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSig := j.sign(signingInput)
	providedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid signature encoding")
	}

	if !hmac.Equal(expectedSig, providedSig) {
		return nil, fmt.Errorf("invalid signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid claims encoding")
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("invalid claims json")
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

func (j *JWTHandler) ParseClaims(token string) (*JWTClaims, error) {
	return j.Verify(token)
}

func (j *JWTHandler) sign(input string) []byte {
	mac := hmac.New(sha256.New, j.secret)
	mac.Write([]byte(input))
	return mac.Sum(nil)
}

type APIKeyValidator struct {
	store *APIKeyStore
}

func NewAPIKeyValidator(store *APIKeyStore) *APIKeyValidator {
	return &APIKeyValidator{store: store}
}

func (v *APIKeyValidator) Validate(key string) (*APIKey, error) {
	apiKey, err := v.store.Get(key)
	if err != nil {
		return nil, err
	}
	if apiKey == nil {
		return nil, fmt.Errorf("invalid api key")
	}
	if !apiKey.IsActive {
		return nil, fmt.Errorf("api key is inactive")
	}
	return apiKey, nil
}

func (v *APIKeyValidator) Create(service string) (*APIKey, error) {
	if service == "" {
		return nil, fmt.Errorf("service name is required")
	}
	key := &APIKey{
		Key:       generateKey(),
		Service:   service,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	if err := v.store.Store(key); err != nil {
		return nil, err
	}
	return key, nil
}

func generateKey() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(time.Now().UnixNano() % 256)
		time.Sleep(1)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}
