package mtls

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrCertNotFound = errors.New("certificate not found")
	ErrCertRevoked  = errors.New("certificate is revoked")
)

type CertManager struct {
	ca *CertificateAuthority
}

func NewCertManager(ca *CertificateAuthority) *CertManager {
	return &CertManager{ca: ca}
}

func (m *CertManager) IssueCert(ctx context.Context, serviceName, cn, org string, validityDays int, isServer bool) (*Certificate, error) {
	if isServer {
		return m.ca.IssueServerCert(serviceName, cn, org, validityDays)
	}
	return m.ca.IssueClientCert(serviceName, cn, org, validityDays)
}

func (m *CertManager) RenewCert(ctx context.Context, certID string, validityDays int) (*Certificate, error) {
	oldCert, ok := m.ca.IssuedCerts[certID]
	if !ok {
		return nil, ErrCertNotFound
	}
	return m.IssueCert(ctx, oldCert.ServiceName, oldCert.CommonName, oldCert.Organization, validityDays, true)
}

func (m *CertManager) RevokeCert(ctx context.Context, certID string) error {
	return m.ca.Revoke(certID)
}

func (m *CertManager) VerifyCert(ctx context.Context, certID string) error {
	cert, ok := m.ca.IssuedCerts[certID]
	if !ok {
		return ErrCertNotFound
	}
	if cert.IsRevoked {
		return ErrCertRevoked
	}
	return nil
}

func (m *CertManager) ListCerts(ctx context.Context) []*Certificate {
	return m.ca.ListCerts()
}

func (m *CertManager) VerifyChain(ctx context.Context, certPEM []byte) error {
	if err := m.ca.VerifyChain(certPEM); err != nil {
		return fmt.Errorf("chain verification failed: %w", err)
	}
	return nil
}
