# Security Rules & Threat Mitigations

## 1. Authentication & Session Security
- **JWT Standard**: Tokens must contain short lifetimes (15 minutes access token, 7 days refresh token).
- **Token Storage**: Store access tokens in short-lived memory; refresh tokens must be stored in secure `HTTP-Only`, `SameSite=Strict`, `Secure` cookies.
- **JWKS Key Rotation**: Verify JWT signatures using a JWKS (JSON Web Key Set) endpoint to allow dynamic key rotation without service downtime.

## 2. API Injection Defense
- **SQL Injection**: String concatenation in raw SQL queries is strictly prohibited. Parameterized queries or Query DSLs must be used.
  ```go
  // FORBIDDEN: db.Raw("SELECT * FROM products WHERE id = '" + input + "'")
  // APPROVED:
  db.Raw("SELECT * FROM products WHERE id = ?", inputID)
  ```
- **XSS & CSRF Mitigation**:
  - Sanitize all user-inputted rich text (e.g., product review comments) using DOMPurify on the frontend and HTML sanitization libraries on the backend.
  - Use double-submit CSRF cookie patterns or standard Next.js Middleware header check for state-changing requests.

## 3. Rate-Limiting & API Shielding
- **Redis Sliding-Window Rate Limiter**: Rate limit critical endpoints (`/api/v1/auth/login`, `/api/v1/checkout`) based on IP and authenticated `UserID`.
- Limit login requests to maximum 5 attempts per 5 minutes per IP.
- Limit checkout requests to maximum 1 attempt per 5 seconds per User.
