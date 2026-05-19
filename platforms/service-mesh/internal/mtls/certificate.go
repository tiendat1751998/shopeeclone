package mtls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type keyPair struct {
	CertPEM []byte
	KeyPEM  []byte
}

type CertificateAuthority struct {
	RootCert       *x509.Certificate
	RootKey        *rsa.PrivateKey
	RootPEM        []byte
	IssuedCerts    map[string]*Certificate
	revokedCerts   map[string]bool
}

func NewCertificateAuthority(org, cn string, validityDays int) (*CertificateAuthority, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate CA key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial: %w", err)
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{org},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, validityDays),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create CA cert: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA cert: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	serialStr := serial.Text(16)
	fingerprint := sha256Hex(certDER)

	caCert := &Certificate{
		ID:           "ca-root",
		ServiceName:  "ca",
		CommonName:   cn,
		Organization: org,
		ValidityDays: validityDays,
		IssuedAt:     now,
		ExpiresAt:    template.NotAfter,
		Serial:       serialStr,
		Fingerprint:  fingerprint,
		IsRevoked:    false,
		IsCA:         true,
	}

	ca := &CertificateAuthority{
		RootCert:     cert,
		RootKey:      key,
		RootPEM:      certPEM,
		IssuedCerts:  map[string]*Certificate{"ca-root": caCert},
		revokedCerts: make(map[string]bool),
	}
	_ = keyPEM
	return ca, nil
}

func (ca *CertificateAuthority) IssueClientCert(serviceName, cn, org string, validityDays int) (*Certificate, error) {
	return ca.issueCert(serviceName, cn, org, validityDays, x509.ExtKeyUsageClientAuth)
}

func (ca *CertificateAuthority) IssueServerCert(serviceName, cn, org string, validityDays int) (*Certificate, error) {
	return ca.issueCert(serviceName, cn, org, validityDays, x509.ExtKeyUsageServerAuth)
}

func (ca *CertificateAuthority) issueCert(serviceName, cn, org string, validityDays int, extKeyUsage x509.ExtKeyUsage) (*Certificate, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial: %w", err)
	}

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{org},
		},
		NotBefore:             now,
		NotAfter:              now.AddDate(0, 0, validityDays),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{extKeyUsage},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, ca.RootCert, &key.PublicKey, ca.RootKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cert: %w", err)
	}

	serialStr := serial.Text(16)
	fingerprint := sha256Hex(certDER)

	certObj := &Certificate{
		ID:           fmt.Sprintf("cert-%s-%s", serviceName, serialStr),
		ServiceName:  serviceName,
		CommonName:   cn,
		Organization: org,
		ValidityDays: validityDays,
		IssuedAt:     now,
		ExpiresAt:    cert.NotAfter,
		Serial:       serialStr,
		Fingerprint:  fingerprint,
		IsRevoked:    false,
		IsCA:         false,
	}

	ca.IssuedCerts[certObj.ID] = certObj
	_ = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	_ = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})

	return certObj, nil
}

func (ca *CertificateAuthority) VerifyChain(certPEM []byte) error {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse cert: %w", err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca.RootCert)

	opts := x509.VerifyOptions{
		Roots: roots,
	}
	_, err = cert.Verify(opts)
	return err
}

func (ca *CertificateAuthority) Revoke(certID string) error {
	cert, ok := ca.IssuedCerts[certID]
	if !ok {
		return fmt.Errorf("certificate not found: %s", certID)
	}
	cert.IsRevoked = true
	ca.revokedCerts[certID] = true
	return nil
}

func (ca *CertificateAuthority) ListCerts() []*Certificate {
	var result []*Certificate
	for _, c := range ca.IssuedCerts {
		result = append(result, c)
	}
	return result
}

func (ca *CertificateAuthority) GetRevokedList() []string {
	var result []string
	for id := range ca.revokedCerts {
		result = append(result, id)
	}
	return result
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
